package indexer

import (
	"context"
	"log"
	"time"

	"solana-token-indexer/internal/db"
	"solana-token-indexer/internal/util"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func FetchAndStoreTransactions(mint string) {
	rpcClient := rpc.New(rpc.MainNetBeta_RPC) 

	mintPK, err := solana.PublicKeyFromBase58(mint)
	if err != nil {
		log.Printf("[indexer] Invalid mint public key %s: %v", mint, err)
		return
	}

	accountsResult, err := rpcClient.GetTokenLargestAccounts(
		context.Background(),
		mintPK,
		rpc.CommitmentFinalized, // Pass commitment directly
	)
	if err != nil {
		log.Printf("[rpc] GetTokenLargestAccounts for mint %s failed: %v", mint, err)
		return
	}

	for _, acctPair := range accountsResult.Value {
		sigsResult, err := rpcClient.GetSignaturesForAddress(
			context.Background(),
			acctPair.Address, 
		)
		if err != nil {
			log.Printf("[rpc] GetSignaturesForAddress for account %s failed: %v", acctPair.Address.String(), err)
			continue
		}

		for _, sigInfo := range sigsResult {
			txResult, err := rpcClient.GetTransaction(
				context.Background(),
				sigInfo.Signature, 
				&rpc.GetTransactionOpts{
					Encoding:   solana.EncodingBase64, 
					Commitment: rpc.CommitmentFinalized,
				},
			)
			if err != nil || txResult == nil || txResult.Transaction == nil { 
				log.Printf("[rpc] GetTransaction for sig %s failed or tx/tx.Transaction is nil: %v", sigInfo.Signature.String(), err)
				continue
			}

			decodedTx, err := txResult.Transaction.GetTransaction()
			if err != nil {
				log.Printf("[rpc] Failed to decode transaction for sig %s: %v", sigInfo.Signature.String(), err)
				continue
			}
			if decodedTx == nil { 
				log.Printf("[rpc] Decoded transaction is nil for sig %s", sigInfo.Signature.String())
				continue
			}

			for _, inst := range decodedTx.Message.Instructions {
				programPK := decodedTx.Message.AccountKeys[inst.ProgramIDIndex]
				if !programPK.Equals(solana.TokenProgramID) {
					continue
				}

				amount, err := util.ParseSPLInstruction(string(inst.Data))
				if err != nil {
					continue
				}

				if len(inst.Accounts) < 2 { 
					continue
				}
				srcPK := decodedTx.Message.AccountKeys[inst.Accounts[0]]
				dstPK := decodedTx.Message.AccountKeys[inst.Accounts[1]]

				var timeStr string
				if txResult.BlockTime != nil {
					timeStr = time.Unix(int64(*txResult.BlockTime), 0).UTC().Format(time.RFC3339)
				}

				db.InsertTransaction(sigInfo.Signature.String(), txResult.Slot, timeStr, srcPK.String(), dstPK.String(), amount, mint)
			}
		}
	}
	log.Printf("[indexer] Finished indexing for mint %s\n", mint)
}

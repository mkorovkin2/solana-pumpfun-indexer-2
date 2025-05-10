# Solana Token Indexer

This Go service monitors and indexes historical SPL token transfers for a given token mint address. It interacts directly with the Solana blockchain using the RPC interface and stores decoded transactions in a PostgreSQL database.

---

## Features

- Token transaction history retrieval via `/query`
- Real SPL token instruction decoding (`Transfer`, `TransferChecked`)
- PostgreSQL persistence
- Modular, production-grade codebase

---

## Setup

### Requirements

- Go 1.20+
- PostgreSQL (running locally)
- Internet or custom Solana RPC node

### Run

```bash
git clone https://github.com/mkorovkin2/solana-pump-fun-indexer-2.git
cd solana-pump-fun-indexer-2
createdb solana
go mod tidy
go run cmd/indexer/main.go
```

### Sample API call
```ash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"mint": "So11111111111111111111111111111111111111112"}'
 ```

### PostgreSQL schema

Assume we have created database `solana_index`. You can then creaete a table there as such:
```
CREATE TABLE IF NOT EXISTS token_transactions (
  signature TEXT PRIMARY KEY,
  slot BIGINT,
  block_time TIMESTAMP,
  source TEXT,
  destination TEXT,
  amount BIGINT,
  token_mint TEXT
);
```

### To get your RPC endpoint

* set one on [Helius](https://www.helius.dev/)
* then set your `SOLANA_RPC_URL` environment variable to the endpoint you get from Helius

### Testing

* run `./test_local_run.sh`


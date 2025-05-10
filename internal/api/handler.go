package api

import (
    "encoding/json"
    "log"
    "net/http"

    "solana-token-indexer/internal/indexer"
)

type QueryRequest struct {
    Mint string `json:"mint"`
}

func HandleQuery(w http.ResponseWriter, r *http.Request) {
    var req QueryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Mint == "" {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    log.Printf("[api] Received query for token: %s\n", req.Mint)
    go indexer.FetchAndStoreTransactions(req.Mint)
    w.Write([]byte("Started indexing transactions..."))
}

package main

import (
    "log"
    "net/http"

    "solana-token-indexer/internal/api"
    "solana-token-indexer/internal/db"
    "github.com/gorilla/mux"
)

func main() {
    db.InitDB()
    r := mux.NewRouter()
    r.HandleFunc("/query", api.HandleQuery).Methods("POST")

    log.Println("[server] Listening on :8080...")
    http.ListenAndServe(":8080", r)
}

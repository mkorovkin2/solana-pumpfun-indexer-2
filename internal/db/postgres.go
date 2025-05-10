package db

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("postgres", "postgres://postgres:devpassword@localhost:5432/solana_index?sslmode=disable")
    if err != nil {
        log.Fatalf("[db] connection error: %v", err)
    }
    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS token_transactions (
        signature TEXT PRIMARY KEY,
        slot BIGINT,
        block_time TIMESTAMP,
        source TEXT,
        destination TEXT,
        amount BIGINT,
        token_mint TEXT
    )`)
    if err != nil {
        log.Fatalf("[db] table creation error: %v", err)
    }
    log.Println("[db] Connected and initialized.")
}

func InsertTransaction(signature string, slot uint64, blockTime string, source, dest string, amount uint64, mint string) {
    _, err := DB.Exec(`INSERT INTO token_transactions (signature, slot, block_time, source, destination, amount, token_mint)
        VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (signature) DO NOTHING`,
        signature, slot, blockTime, source, dest, amount, mint)
    if err != nil {
        log.Printf("[db] insertion error: %v", err)
    }
}

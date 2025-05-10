#!/bin/bash

# Configuration
DB_NAME="solana_index"
TEST_MINT="So11111111111111111111111111111111111111112" # Solana (SOL) mint address
# TEST_MINT="EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzL7DKnvLKeekH" # Example for a smaller token (USDC on devnet, adjust if needed for mainnet)
SERVER_LOG="server.log"
INDEXER_PID=""
EXPECTED_TRANSACTIONS=1 # Adjust based on what you expect for the TEST_MINT after a short run

# Function to clean up server
cleanup() {
    echo "Cleaning up..."
    if [ ! -z "$INDEXER_PID" ]; then
        echo "Stopping indexer server (PID: $INDEXER_PID)..."
        kill $INDEXER_PID
        # Wait a bit for the server to shut down
        sleep 2
    fi
    rm -f $SERVER_LOG
    echo "Cleanup done."
}

# Trap EXIT signal to ensure cleanup
trap cleanup EXIT

echo "Starting test script..."

# 1. Setup Database (Optional: Clear previous test data for this mint)
echo "Ensuring database table exists..."
psql -d $DB_NAME -c "CREATE TABLE IF NOT EXISTS token_transactions (signature TEXT PRIMARY KEY, slot BIGINT, block_time TIMESTAMP, source TEXT, destination TEXT, amount BIGINT, token_mint TEXT);"
if [ $? -ne 0 ]; then
    echo "Failed to ensure database table exists. Exiting."
    exit 1
fi

echo "Optional: Clearing previous test data for mint $TEST_MINT..."
psql -d $DB_NAME -c "DELETE FROM token_transactions WHERE token_mint = '$TEST_MINT';"
if [ $? -ne 0 ]; then
    echo "Warning: Failed to clear previous test data. Continuing..."
fi

# 2. Run Indexer Server in the background
echo "Starting indexer server in background..."
go run cmd/indexer/main.go > $SERVER_LOG 2>&1 &
INDEXER_PID=$!
echo "Server PID: $INDEXER_PID. Log: $SERVER_LOG"

# Wait a few seconds for the server to initialize
echo "Waiting for server to start..."
sleep 5 

# Check if server started successfully (basic check, looks for listening message)
if ! grep -q "Listening on :8080" $SERVER_LOG; then
    echo "Server did not start successfully or quickly enough. Check $SERVER_LOG. Exiting."
    exit 1
fi
echo "Server started."

# 3. Trigger Indexing via API call
echo "Triggering indexing for mint: $TEST_MINT..."
RESPONSE=$(curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d "{\"mint\": \"$TEST_MINT\"}")

echo "API Response: $RESPONSE"
if [[ "$RESPONSE" != "Started indexing transactions..." ]]; then
    echo "Unexpected API response. Expected 'Started indexing transactions...'. Exiting."
    exit 1
fi

# 4. Wait for indexing to occur
# This is an estimate. For a busy token, more time might be needed.
# For a test, you might want to use a mint with few transactions.
INDEXING_WAIT_TIME=30 # seconds
echo "Waiting $INDEXING_WAIT_TIME seconds for indexing to fetch some transactions..."
sleep $INDEXING_WAIT_TIME

# 5. Check Database for Indexed Transactions
echo "Checking database for transactions for mint $TEST_MINT..."
COUNT=$(psql -d $DB_NAME -t -c "SELECT COUNT(*) FROM token_transactions WHERE token_mint = '$TEST_MINT';")

# Trim whitespace from COUNT
COUNT=$(echo $COUNT | xargs)

echo "Found $COUNT transactions for mint $TEST_MINT."

if [ -z "$COUNT" ] || ! [[ "$COUNT" =~ ^[0-9]+$ ]]; then
    echo "Error: Could not retrieve a valid count from the database."
    exit 1
fi

if [ "$COUNT" -ge "$EXPECTED_TRANSACTIONS" ]; then
    echo "SUCCESS: Found $COUNT transactions (expected at least $EXPECTED_TRANSACTIONS)."
    exit 0
else
    echo "FAILURE: Found only $COUNT transactions (expected at least $EXPECTED_TRANSACTIONS)."
    echo "Check server log ($SERVER_LOG) and database."
    exit 1
fi

# Cleanup will be called automatically on exit

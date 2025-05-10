#!/bin/bash

# Sample usage: ./query_server.sh 41fZewbrb8x24yE9KeMJExoVXCDggebPMEwRVvX8pump

# Check if a mint address is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <mint_address>"
  echo "Example: $0 So11111111111111111111111111111111111111112"
  exit 1
fi

MINT_ADDRESS="$1"
SERVER_URL="http://localhost:8080/query"

echo "Sending query to $SERVER_URL for mint: $MINT_ADDRESS"

curl -X POST "$SERVER_URL" \
  -H "Content-Type: application/json" \
  -d "{\"mint\": \"$MINT_ADDRESS\"}"

# Add a newline for cleaner terminal output after curl
echo 

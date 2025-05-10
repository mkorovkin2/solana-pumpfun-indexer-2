# Use Go base image
FROM golang:1.20

# Set working directory
WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the source
COPY . .

# Build the app
RUN go build -o indexer ./cmd/indexer

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./indexer"]

# Use an official Golang image for building the application
FROM golang:1.18 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy project files
COPY . .

# Download dependencies and build the CLI wallet
RUN go mod tidy && go build -o wallet main.go

# Create a minimal runtime container
FROM alpine:latest

# Set environment variables for security & Ethereum RPC
ENV WALLET_ENCRYPTION_KEY=""
ENV ETH_RPC_URL="https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"
ENV ETHERSCAN_API_KEY="your_etherscan_api_key_here"

# Set working directory inside the runtime container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/wallet .

# Make sure the binary is executable
RUN chmod +x wallet

# Run the CLI tool by default (displays help if no command is passed)
ENTRYPOINT ["/app/wallet"]

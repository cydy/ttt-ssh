# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tictactoe-ssh .

# Final stage - scratch base for minimal container
FROM scratch

# Copy ca-certificates from alpine
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder
COPY --from=builder /app/tictactoe-ssh /tictactoe-ssh

# Expose SSH port
EXPOSE 2222

# Run the application
CMD ["/tictactoe-ssh"]

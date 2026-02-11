# Build the Go Binary
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY . .
# Build specifically for Linux
# The build tags will automatically select 'linux.go' and ignore 'mock.go'
RUN CGO_ENABLED=0 GOOS=linux go build -o aar-daemon main.go

# Final Runtime Image
FROM alpine:latest

# Install system tools required for the router to function
# iproute2: for 'ip' command
# wireguard-tools: for 'wg-quick'
# iw: for Wi-Fi management
# iputils: for 'ping'
RUN apk add --no-cache iproute2 wireguard-tools iw iputils

WORKDIR /root/
COPY --from=builder /app/aar-daemon .

# Expose API
EXPOSE 8080

# Run the daemon
CMD ["./aar-daemon"]

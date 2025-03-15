FROM golang:1.24-alpine AS go-builder

WORKDIR /app

# Copy go.mod and go.sum (relative to the Docker build context)
COPY ./go.mod ./go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the Go project
COPY ./ ./

# Build the Gin service from cmd/main.go
RUN go build -o gin-service ./cmd/main.go

# Minimal final image
FROM alpine:latest

WORKDIR /root/

# Copy the built binary
COPY --from=go-builder /app/gin-service .

# Set environment variable for default port
ENV PORT=8081

EXPOSE 8081

CMD ["./gin-service"]
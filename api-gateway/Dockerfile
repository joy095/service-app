FROM golang:1.24-alpine AS builder

WORKDIR /app/api-gateway

# Copy go.mod and go.sum to cache dependencies
COPY ./api-gateway/go.mod ./api-gateway/go.sum ./
RUN go mod download

# Copy the rest of the code
COPY /api-gateway .

# Build the Go application
RUN go build -o main cmd/main.go

EXPOSE 8080

# Run the application
CMD ["./main"]
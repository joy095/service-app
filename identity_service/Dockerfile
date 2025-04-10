FROM golang:1.24-alpine AS builder

WORKDIR /app/identity_service

# Copy go.mod and go.sum to cache dependencies
COPY ./identity_service/go.mod ./identity_service/go.sum ./
RUN go mod download

# Copy the rest of the code
COPY /identity_service .

# Build the Go application
RUN go build -o main cmd/main.go

EXPOSE 8081

# Run the application
CMD ["./main"]
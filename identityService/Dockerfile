FROM golang:1.24-alpine AS builder

WORKDIR /app/identityService

# Copy go.mod and go.sum to cache dependencies
COPY ./identityService/go.mod ./identityService/go.sum ./
RUN go mod download

# Copy the rest of the code
COPY /identityService .

# Build the Go application
RUN go build -o main cmd/main.go

EXPOSE 8081

# Run the application
CMD ["./main"]
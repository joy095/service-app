FROM golang:1.24-alpine AS builder

WORKDIR /

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
FROM golang:1.24-alpine AS builder

WORKDIR /app/word_filter_service

# Copy go.mod and go.sum to cache dependencies
COPY ./word_filter_service/go.mod ./word_filter_service/go.sum ./
RUN go mod download

# Copy the rest of the code
COPY /word_filter_service .

# Build the Go application
RUN go build -o main cmd/main.go

EXPOSE 8082

# Run the application
CMD ["./main"]
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./message-service/go.mod ./message-service/go.sum ./
RUN go mod download

COPY ./message-service .

RUN go build -o main ./cmd/main.go

# Final minimal image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8085

CMD ["./main"]

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ./api-gateway/go.mod ./api-gateway/go.sum ./
RUN go mod download

COPY ./api-gateway .

RUN go build -o main ./cmd/main.go

# Final minimal image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

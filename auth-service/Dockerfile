FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-service ./cmd/server

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/auth-service .
COPY .env .env

EXPOSE 8080

CMD ["./auth-service"]

FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o subscribe-service ./cmd/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/subscribe-service .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./subscribe-service"]
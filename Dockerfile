FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bot .
RUN mkdir -p config
COPY config/.env.example config/.env
RUN adduser -D -g '' botuser
USER botuser

CMD ["./bot"] 
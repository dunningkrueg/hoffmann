FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bot .
RUN mkdir -p config
RUN adduser -D -g '' botuser
USER botuser

CMD ["./bot"] 
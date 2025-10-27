# Build stage
FROM golang:1.25-alpine as builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY . .

RUN go mod tidy && go mod download

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bot-service ./cmd/api

# Final Stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/bot-service /bot-service

ENTRYPOINT ["/bot-service"]
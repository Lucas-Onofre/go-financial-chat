# Build stage
FROM golang:1.24-alphine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLE=0 go build -ldflags="-s -w" -o chat-service ./cmd/api/main.go

# Final Stage
FROM scratch

COPY --from=builder /app/chat-service /chat-service

ENTRYPOINT ["/chat-service"]
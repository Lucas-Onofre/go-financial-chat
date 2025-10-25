# Build stage
FROM golang:1.25-alpine as builder

WORKDIR /app

COPY . .
RUN go mod tidy && go mod download

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bot-service ./cmd/api

# Final Stage
FROM scratch

COPY --from=builder /app/bot-service /bot-service

ENTRYPOINT ["/bot-service"]
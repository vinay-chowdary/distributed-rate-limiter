# Stage 1: build
FROM golang:1.27 AS builder
WORKDIR /app
COPY . .
RUN go build -o rate-limiter ./cmd/server

# Stage 2: run
FROM debian:bookworm-slim
COPY --from=builder /app/rate-limiter .
CMD ["./rate-limiter"]
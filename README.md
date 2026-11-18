# Distributed Rate Limiter

A distributed rate limiter in Go. It uses the **Token Bucket** algorithm, stores state in **Redis**, and serves checks over **gRPC**.

## How it works

- Each client has a bucket with `capacity` tokens. The bucket refills at `refillRate` tokens per second.
- A Lua script on Redis does the refill and the token take as one atomic step, so many server instances can share the same limits without race conditions.
- Keys per client: `ratelimit:<client_id>:tokens` and `ratelimit:<client_id>:last_refill`.
- Defaults: capacity `10`, refill `0.1/s`, which is 1 token every 10s ([cmd/server/main.go](cmd/server/main.go)).

```
client ──gRPC──▶ rate-limiter server ──EVAL (Lua)──▶ Redis
```

## Project layout

| Path | Purpose |
|---|---|
| `cmd/server` | gRPC server, health check, Prometheus metrics, graceful shutdown |
| `cmd/client` | Sample client that sends 12 requests as `user-1` |
| `internal/ratelimiter` | Redis limiter (`redis_limiter.go`) and an in-memory limiter (`limiter.go`) |
| `internal/server` | gRPC handler |
| `proto/ratelimiter` | Protobuf definition and generated code |

## API

```proto
rpc CheckRateLimit(RateLimitCheckRequest{ client_id }) 
    returns (RateLimitCheckResponse{ allowed, remaining_tokens, retry_after_ms });
```

## Run

```bash
# Option 1: Docker (Redis + server)
docker compose up --build

# Option 2: local
docker run -p 6379:6379 redis:7
go run ./cmd/server          # REDIS_ADDR defaults to localhost:6379

# Try it
go run ./cmd/client

# Tests
go test ./...
```

| Port | Endpoint |
|---|---|
| `50051` | gRPC (`RateLimitService` + `grpc.health.v1`) |
| `9090` | Prometheus `/metrics` |

| Env var | Default |
|---|---|
| `REDIS_ADDR` | `localhost:6379` |

## Roadmap (priority order)

1. **TTL on Redis keys**: stops a memory leak from idle clients.
2. **Fail-open policy**: keeps serving when Redis is down.
3. **Config from environment variables**: no hardcoded capacity, rate or ports.
4. **EVALSHA script caching**: sends less data to Redis per call.
5. **Redis health check goroutine**: finds Redis problems before requests fail.
6. **Per-user rate limit tiers**: different limits for each plan.
7. **Token leasing**: fewer Redis round trips, because servers take tokens in batches.
8. **Redis Cluster/Sentinel**: removes the single point of failure.
9. **Custom Prometheus metrics**: counters for allowed and rejected requests.
10. **Multiple rate limit windows**: for example 100/min and 1000/hour.

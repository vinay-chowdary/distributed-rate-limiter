package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const tokenBucketScript = `
		local stored_tokens = redis.call("GET", KEYS[1])
		local last_refill_time = redis.call("GET", KEYS[2])
		local capacity = tonumber(ARGV[1])
		local refil_rate = tonumber(ARGV[2])
		local current_time = tonumber(ARGV[3])

		if stored_tokens == false and last_refill_time == false then
			stored_tokens = capacity
			last_refill_time = current_time
		else
			stored_tokens = tonumber(stored_tokens)
			last_refill_time = tonumber(last_refill_time)
			local elapsed = current_time - last_refill_time
			local gained = elapsed * refil_rate + stored_tokens
			if gained > capacity then
				stored_tokens = capacity
			else
				stored_tokens = gained
			end
			last_refill_time = current_time
		end



		redis.call("SET", KEYS[2], last_refill_time)


		if stored_tokens >= 1 then
			stored_tokens = stored_tokens - 1
			redis.call("SET", KEYS[1], stored_tokens)
			return {1, math.floor(stored_tokens), 0}
		else
			redis.call("SET", KEYS[1], stored_tokens)
			return {0, 0, math.floor((1/refil_rate)*1000)}
		end
`

type RedisRateLimiter struct {
	client     *redis.Client
	capacity   float64
	refillRate float64
}

func NewRedisRateLimiter(client *redis.Client, capacity float64, refillRate float64) *RedisRateLimiter {
	return &RedisRateLimiter{client: client, capacity: capacity, refillRate: refillRate}
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, clientID string) (bool, int64, int64, error) {
	lastRefillKey := fmt.Sprintf("ratelimit:%s:last_refill", clientID)
	tokenKey := fmt.Sprintf("ratelimit:%s:tokens", clientID)
	currentTime := float64(time.Now().UnixMicro()) / 1e6

	result, err := rl.client.Eval(ctx, tokenBucketScript,
		[]string{tokenKey, lastRefillKey},
		rl.capacity, rl.refillRate, currentTime,
	).Slice()

	if err != nil {
		return false, 0, 0, err
	}

	allowed := result[0].(int64)
	remainingTokens := result[1].(int64)
	retryAfter := result[2].(int64)

	return allowed == 1, remainingTokens, retryAfter, nil
}

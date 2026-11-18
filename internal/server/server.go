package server

import (
	"context"

	"github.com/vinay-chowdary/distributed-rate-limiter/internal/ratelimiter"
	pr "github.com/vinay-chowdary/distributed-rate-limiter/proto/ratelimiter"
)

type rateLimitServer struct {
	limiter *ratelimiter.RedisRateLimiter
	pr.UnimplementedRateLimitServiceServer
}

func NewRateLimitServer(limiter *ratelimiter.RedisRateLimiter) *rateLimitServer {
	return &rateLimitServer{limiter: limiter}
}

func (rl *rateLimitServer) CheckRateLimit(ctx context.Context, req *pr.RateLimitCheckRequest) (*pr.RateLimitCheckResponse, error) {
	isAllowed, remainingTokens, retryAfter, err := rl.limiter.Allow(ctx, req.ClientId)

	if err != nil {
		return nil, err
	}

	return &pr.RateLimitCheckResponse{Allowed: isAllowed, RemainingTokens: int32(remainingTokens), RetryAfterMs: retryAfter}, nil
}

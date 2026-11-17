package ratelimiter

import (
	"math"
	"time"
)

type TokenBucket struct {
	capacity       float64
	refillRate     float64
	lastRefillTime time.Time
	storedTokens   float64
}

func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	bucket := TokenBucket{capacity: capacity, refillRate: refillRate, lastRefillTime: time.Now(), storedTokens: capacity}

	return &bucket
}

func (tb *TokenBucket) Allow() bool {
	elapsed := time.Since(tb.lastRefillTime)

	newTokens := elapsed.Seconds() * tb.refillRate

	currentTokens := math.Min(tb.storedTokens+newTokens, tb.capacity)

	tb.storedTokens = currentTokens
	tb.lastRefillTime = time.Now()

	if currentTokens >= 1 {
		tb.storedTokens--
		return true
	}

	return false
}

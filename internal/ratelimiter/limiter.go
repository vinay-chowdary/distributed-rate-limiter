package ratelimiter

import (
	"math"
	"sync"
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

type RateLimiter struct {
	userBuckets       map[string]*TokenBucket
	defaultCapacity   float64
	defaultRefillrate float64
	lock              sync.Mutex
}

func NewRateLimiter(capacity float64, refillRate float64) *RateLimiter {
	return &RateLimiter{defaultCapacity: capacity, defaultRefillrate: refillRate, userBuckets: make(map[string]*TokenBucket)}
}

func (rl *RateLimiter) Allow(clientID string) bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	bucket, exists := rl.userBuckets[clientID]

	if !exists {
		newBucket := NewTokenBucket(rl.defaultCapacity, rl.defaultRefillrate)
		rl.userBuckets[clientID] = newBucket

		return newBucket.Allow()
	}

	return bucket.Allow()

}

package main

import (
	"fmt"
	"time"

	"github.com/vinay-chowdary/distributed-rate-limiter/internal/ratelimiter"
)

func main() {

	bucket := ratelimiter.NewTokenBucket(5, 1)

	for req := range 7 {
		fmt.Printf("Request %d: %v\n", req, bucket.Allow())

	}
	time.Sleep(2 * time.Second)

	for req := range 3 {
		fmt.Printf("Request %d: %v\n", req, bucket.Allow())

	}
}

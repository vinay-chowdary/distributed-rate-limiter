package main

import (
	"fmt"

	"github.com/vinay-chowdary/distributed-rate-limiter/internal/ratelimiter"
)

func main() {

	limiter := ratelimiter.NewRateLimiter(3, 1)

	for req := range 5 {
		isAllowed := limiter.Allow("user-1")

		fmt.Printf("Request user-1 %d: %v\n", req, isAllowed)

	}
	// time.Sleep(2 * time.Second)

	for req := range 2 {
		isAllowed := limiter.Allow("user-2")
		fmt.Printf("Request user-2 %d: %v\n", req, isAllowed)

	}
}

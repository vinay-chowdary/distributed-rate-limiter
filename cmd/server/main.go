package main

import (
	"fmt"
	"sync"

	"github.com/vinay-chowdary/distributed-rate-limiter/internal/ratelimiter"
)

func main() {

	limiter := ratelimiter.NewRateLimiter(100, 10)

	wg := sync.WaitGroup{}

	for routine := range 50 {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			for i := range 10 {
				isAllowed := limiter.Allow("user-1")

				fmt.Printf("user-1, Routine %d, req  %d, allowed %v\n", r, i, isAllowed)
			}
		}(routine)
	}

	wg.Wait()

	// for req := range 5 {
	// 	isAllowed := limiter.Allow("user-1")

	// 	fmt.Printf("Request user-1 %d: %v\n", req, isAllowed)

	// }
	// time.Sleep(2 * time.Second)

	// for req := range 2 {
	// 	isAllowed := limiter.Allow("user-2")
	// 	fmt.Printf("Request user-2 %d: %v\n", req, isAllowed)

	// }
}

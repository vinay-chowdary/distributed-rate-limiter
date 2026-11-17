package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/vinay-chowdary/distributed-rate-limiter/internal/ratelimiter"
)

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatal(err)
	}

	limiter := ratelimiter.NewRedisRateLimiter(rdb, 5, 1)

	wg := sync.WaitGroup{}

	for routine := range 50 {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			for i := range 10 {
				isAllowed, err := limiter.Allow(ctx, "user-1")

				if err != nil {
					fmt.Printf("something went wrong %v", err)
				}

				fmt.Printf("user-1, Routine %d, req  %d, allowed %v\n", r, i, isAllowed)
			}
		}(routine)
	}

	wg.Wait()

	// for req := range 7 {
	// 	isAllowed, err := limiter.Allow(ctx, "user-1")

	// 	if err != nil {
	// 		fmt.Printf("something went wrong %v", err)
	// 	}

	// 	fmt.Printf("Request user-1 %d: %v\n", req, isAllowed)

	// }
	// time.Sleep(3 * time.Second)

	// for req := range 3 {
	// 	isAllowed, err := limiter.Allow(ctx, "user-1")
	// 	if err != nil {
	// 		fmt.Printf("something went wrong %v", err)
	// 	}
	// 	fmt.Printf("Request user-1 %d: %v\n", req, isAllowed)

	// }
}

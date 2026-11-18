package main

import (
	"context"
	"log"
	"net"

	"github.com/redis/go-redis/v9"
	"github.com/vinay-chowdary/distributed-rate-limiter/internal/ratelimiter"
	"github.com/vinay-chowdary/distributed-rate-limiter/internal/server"
	pr "github.com/vinay-chowdary/distributed-rate-limiter/proto/ratelimiter"
	"google.golang.org/grpc"
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

	limiter := ratelimiter.NewRedisRateLimiter(rdb, 10, 0.1)

	rateLimitServer := server.NewRateLimitServer(limiter)

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pr.RegisterRateLimitServiceServer(grpcServer, rateLimitServer)

	log.Println("gRPC server listening on :50051")

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal(err)
	}

	// wg := sync.WaitGroup{}

	// for routine := range 50 {
	// 	wg.Add(1)
	// 	go func(r int) {
	// 		defer wg.Done()
	// 		for i := range 10 {
	// 			isAllowed, err := limiter.Allow(ctx, "user-1")

	// 			if err != nil {
	// 				fmt.Printf("something went wrong %v", err)
	// 			}

	// 			fmt.Printf("user-1, Routine %d, req  %d, allowed %v\n", r, i, isAllowed)
	// 		}
	// 	}(routine)
	// }

	// wg.Wait()

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

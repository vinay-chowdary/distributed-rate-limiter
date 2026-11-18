package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/vinay-chowdary/distributed-rate-limiter/internal/ratelimiter"
	"github.com/vinay-chowdary/distributed-rate-limiter/internal/server"
	pr "github.com/vinay-chowdary/distributed-rate-limiter/proto/ratelimiter"
	"google.golang.org/grpc"
)

func main() {

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
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

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down gracefully...")
	grpcServer.GracefulStop()

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

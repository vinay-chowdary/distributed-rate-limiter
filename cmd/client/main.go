package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/vinay-chowdary/distributed-rate-limiter/proto/ratelimiter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	con, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	client := pb.NewRateLimitServiceClient(con)

	req := &pb.RateLimitCheckRequest{ClientId: "user-1"}

	for range 12 {
		rateLimiter, err := client.CheckRateLimit(context.Background(), req)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(rateLimiter.Allowed, rateLimiter.RemainingTokens, rateLimiter.RetryAfterMs)

	}

}

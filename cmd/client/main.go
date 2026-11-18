package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	pb "github.com/vinay-chowdary/distributed-rate-limiter/proto/ratelimiter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	con, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		slog.Error("Cannot establish connection to server", "error", err)
		os.Exit(1)
	}

	defer con.Close()

	client := pb.NewRateLimitServiceClient(con)

	req := &pb.RateLimitCheckRequest{ClientId: "user-1"}

	for range 12 {
		rateLimiter, err := client.CheckRateLimit(context.Background(), req)

		if err != nil {
			slog.Error("failed to check rate limit", "error", err)
			os.Exit(1)
		}

		fmt.Println(rateLimiter.Allowed, rateLimiter.RemainingTokens, rateLimiter.RetryAfterMs)

	}

}

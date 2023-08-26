package main

import (
	"context"
	"errors"
	"log"

	"github.com/FluffyKebab/todo/infra/inputport/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("health checked failed: %s", err.Error())
	}
}

func run() error {
	cc, err := grpc.Dial("localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer cc.Close()

	client := pb.NewTodoServiceClient(cc)
	res, err := client.HealthCheck(context.Background(), &pb.HealthCheckRequest{})
	if err != nil {
		return err
	}

	if res.Status == pb.HealthCheckResponse_NOT_SERVING {
		return errors.New("server not serving")
	}

	if res.Status == pb.HealthCheckResponse_UNKNOWN {
		return errors.New("server status is unknown")
	}

	return nil
}

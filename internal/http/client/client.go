package client

import (
	"log"

	"github.com/diegom0ta/go-grpc-gofiber/internal/pb"
	"google.golang.org/grpc"
)

var Client pb.UserServiceClient

func Connect(address string) {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	Client = pb.NewUserServiceClient(conn)
}

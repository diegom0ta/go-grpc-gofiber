package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	db "github.com/diegom0ta/go-grpc-gofiber/internal/database"
	"github.com/diegom0ta/go-grpc-gofiber/internal/models"
	"github.com/diegom0ta/go-grpc-gofiber/internal/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *gorm.DB
}

var (
	port = 1531
	wg   = sync.WaitGroup{}
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db.Connect()

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, &server{db: db.Db})

	reflection.Register(s)

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		sig := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", sig)
		cancel()
		db.Disconnect()
		s.GracefulStop()

		wg.Done()
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	wg.Wait()
	log.Println("clean shutdown")
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	uuid := uuid.New()
	user := models.User{ID: uuid.String(), Name: req.Name, Email: req.Email}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{Id: user.ID}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var user models.User
	if err := s.db.First(&user).Where("id = ?", req.Id).Error; err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{Id: user.ID, Name: user.Name, Email: user.Email}, nil
}

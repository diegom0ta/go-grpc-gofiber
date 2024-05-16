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

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

	s := grpc.NewServer()

	pb.RegisterUserServer(s, &user.Server{})

	reflection.Register(s)

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		sig := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", sig)
		cancel()
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

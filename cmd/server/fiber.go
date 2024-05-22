package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegom0ta/go-grpc-gofiber/internal/http/client"
	h "github.com/diegom0ta/go-grpc-gofiber/internal/http/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	grpcAddress = "localhost:1531"
	port        = 8080
)

func main() {
	client.Connect(grpcAddress)

	app := fiber.New()

	app.Use(cors.New())

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
			log.Printf("Error listening server: %v", err)
		}
	}()

	log.Println("Fiber server is running on port :8080")

	app.Post("/user", h.Register)
	app.Get("/user/:id", h.GetUser)

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Println("Gracefully shutting down started...")

	err := app.Shutdown()
	if err != nil {
		log.Printf("Error while shutting down: %v", err)
	}

	log.Println("Server closed")
}

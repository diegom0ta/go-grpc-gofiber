package main

import (
	"context"
	"log"

	"github.com/diegom0ta/go-grpc-gofiber/internal/pb"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

const (
	grpcAddress = "localhost:1531"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	// Create Fiber app
	app := fiber.New()

	// Define routes
	app.Post("/user", func(c *fiber.Ctx) error {
		req := new(struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		})
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		grpcReq := &pb.CreateUserRequest{
			Name:  req.Name,
			Email: req.Email,
		}
		res, err := client.CreateUser(context.Background(), grpcReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": res.Id})
	})

	app.Get("/user", func(c *fiber.Ctx) error {
		req := new(struct {
			Email string `json:"email"`
		})
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		grpcReq := &pb.GetUserRequest{Email: req.Email}
		res, err := client.GetUser(context.Background(), grpcReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"id":    res.Id,
			"name":  res.Name,
			"email": res.Email,
		})
	})

	// Start Fiber app
	log.Println("Fiber server is running on port :8080")
	log.Fatal(app.Listen(":8080"))
}

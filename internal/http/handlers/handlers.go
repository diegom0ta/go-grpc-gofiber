package handlers

import (
	"context"

	cl "github.com/diegom0ta/go-grpc-gofiber/internal/http/client"
	"github.com/diegom0ta/go-grpc-gofiber/internal/pb"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
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

	res, err := cl.Client.CreateUser(context.Background(), grpcReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": res.Id})
}

func GetUser(c *fiber.Ctx) error {
	// Get id from query
	id := c.Params("id")

	// Get id from body
	// req := new(struct {
	// 	ID string `json:"id"`
	// })
	// if err := c.BodyParser(req); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	// }

	grpcReq := &pb.GetUserRequest{Id: id}
	res, err := cl.Client.GetUser(context.Background(), grpcReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"id":    res.Id,
		"name":  res.Name,
		"email": res.Email,
	})
}

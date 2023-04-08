package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	treblle_fiber "github.com/RafaelPiloto10/treblle-go-fiber/trebble_fiber"
)

func main() {
	// Define a new Fiber app with config.
	app := fiber.New(fiber.Config{})
	godotenv.Load()

	treblle_fiber.Configure(treblle_fiber.Configuration{
		APIKey:    os.Getenv("API_KEY"),
		ProjectID: os.Getenv("PROJECT_ID"),
	})

	app.Use(logger.New(), treblle_fiber.Middleware())
	app.Add("POST", "/ping", Ping)

	// Start server (with or without graceful shutdown).
	app.Listen("localhost:3000")
}

type PingRequest struct {
	count int8
}

func Ping(c *fiber.Ctx) error {
	body := &PingRequest{}
	// Checking received data from JSON body.
	if err := c.BodyParser(body); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map {
		"error": false,
		"msg": fmt.Sprintf("Pinged %v\n", body.count),
	})
}

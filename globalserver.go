package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// startGlobalServer starts a Fiber server with routes accessible globally
func startGlobalServer() {
	app := fiber.New()

	// Define global routes
	app.Get("/global", func(c *fiber.Ctx) error {
		return c.SendString("This endpoint is accessible globally")
	})

	// Listen on all interfaces (0.0.0.0)
	if err := app.Listen(":4000"); err != nil {
		log.Fatalf("Global server error: %v", err)
	}
}

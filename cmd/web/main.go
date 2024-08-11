package main

import (
	"github.com/belovetech/go-ci/internal/web/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	handlers.SetupPipelineRoutes(app)

	app.Listen(":3000")
}

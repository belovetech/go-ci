package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type WithRepoURL struct {
	Url string `json:"url" xml:"url" form:"url"`
}

func SetupPipelineRoutes(app *fiber.App) {
	app.Post("/healthz", postCheckItWorks)
}

func postCheckItWorks(c *fiber.Ctx) error {
	body := &WithRepoURL{}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse body",
		})
	}
	return c.SendString(fmt.Sprintf("Working with repo: %s\n", body.Url))
}

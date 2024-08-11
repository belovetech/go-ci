package handlers

import (
	"fmt"

	"github.com/belovetech/go-ci/internal/ci"
	"github.com/gofiber/fiber/v2"
)

type WithRepoURL struct {
	Url    string `json:"url" xml:"url" form:"url"`
	Branch string `json:"branch" xml:"branch" form:"branch"`
}

func SetupPipelineRoutes(app *fiber.App) {
	pipelinesGroup := app.Group("/api/v1/pipelines")
	pipelinesGroup.Get("/healthz", healthCheck)
	pipelinesGroup.Post("/check-it-works", postCheckItWorks)

}

func postCheckItWorks(c *fiber.Ctx) error {
	body := &WithRepoURL{}

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse body",
		})
	}

	var ws ci.Workspace
	ws, err := ci.NewWorkspaceFromGit("./tmp", body.Url, body.Branch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Working with repo: %s", body.Url),
		"branch":  ws.Branch(),
		"commit":  ws.Commit(),
		"dir":     ws.Dir(),
	})
}

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("Server is up and running")
}

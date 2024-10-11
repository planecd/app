package main

import (
	"io"
	"net/http"

	"github.com/google/go-github/v66/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/planecd/api/clients"
)

// d656e7c43692e7484c68936022b5c583efa94f02

func main() {
	multilog.RegisterLogger(multilog.LogMethod("console"), multilog.NewConsoleLogger(&multilog.NewConsoleLoggerArgs{
		Level:  multilog.DEBUG,
		Format: multilog.FormatText,
		FilterDropPatterns: []*string{
			multilog.PtrString("producer"), // Drop rabbitmq producer logs.
		},
	}))

	gitHubClient := &clients.GitHubClient{
		Secret: "asdfasdfasdf",
	}

	err := gitHubClient.Init(1022275, 55841238, "private-key.pem")
	if err != nil {
		multilog.Fatal("main", "Failed to create GitHub client", map[string]any{
			"error": err,
		})
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		content, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		multilog.Debug("asdf", "asdf", map[string]any{
			"content": string(content),
		})

		return c.JSON(http.StatusOK, string(content))
	})

	e.POST("/github/webhook", func(c echo.Context) error {
		// payload, err := io.ReadAll(c.Request().Body)
		// if err != nil {
		// 	return c.JSON(http.StatusInternalServerError, err.Error())
		// }

		err = gitHubClient.Handle(c.Request(), func(event *github.WorkflowRunEvent) {
			multilog.Debug("Received GitHub webhook", "webhook", map[string]any{
				"event": event,
			})
		})

		return c.JSON(http.StatusOK, "Webhook received")
	})

	e.GET("/github/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Healthy")
	})

	e.Logger.Fatal(e.Start(":11080"))
}

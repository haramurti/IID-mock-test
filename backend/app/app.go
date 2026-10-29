// Package app builds the Fiber API app. It is the one entry point shared by
// cmd/server (single binary / Docker) and api/index.go (Vercel function).
// It lives outside internal/ so the Vercel function at the repo root can import it.
package app

import (
	"fmt"
	_ "time/tzdata" // Asia/Jakarta works even in minimal containers / serverless runtimes

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/haramurti/IID-mock-test/backend/internal/config"
	"github.com/haramurti/IID-mock-test/backend/internal/httpapi"
	"github.com/haramurti/IID-mock-test/backend/internal/leads"
)

// New returns a Fiber app with middleware and all /api routes registered.
func New(cfg config.Config) (*fiber.App, error) {
	source, err := leads.NewSource(cfg.DataSource, cfg.ExpaToken)
	if err != nil {
		return nil, fmt.Errorf("data source: %w", err)
	}

	a := fiber.New(fiber.Config{AppName: "lead-dashboard"})
	a.Use(recover.New())
	a.Use(logger.New())
	httpapi.New(source).Register(a)
	return a, nil
}

// FromEnv loads config from environment variables and builds the app.
func FromEnv() (*fiber.App, config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, cfg, fmt.Errorf("config: %w", err)
	}
	a, err := New(cfg)
	return a, cfg, err
}

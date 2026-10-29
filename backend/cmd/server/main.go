package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/haramurti/IID-mock-test/backend/internal/config"
	"github.com/haramurti/IID-mock-test/backend/internal/httpapi"
	"github.com/haramurti/IID-mock-test/backend/internal/leads"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	app := fiber.New(fiber.Config{AppName: "lead-dashboard"})
	app.Use(recover.New())
	app.Use(logger.New())

	source, err := leads.NewSource(cfg.DataSource, cfg.ExpaToken)
	if err != nil {
		log.Fatalf("data source: %v", err)
	}
	httpapi.New(source).Register(app)

	log.Printf("listening on :%s (data source: %s)", cfg.Port, cfg.DataSource)
	log.Fatal(app.Listen(":" + cfg.Port))
}

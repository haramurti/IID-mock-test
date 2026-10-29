package main

import (
	"log"
	"net/http"
	_ "time/tzdata" // Asia/Jakarta works even in minimal containers without tzdata

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/haramurti/IID-mock-test/backend/internal/config"
	"github.com/haramurti/IID-mock-test/backend/internal/httpapi"
	"github.com/haramurti/IID-mock-test/backend/internal/leads"
	"github.com/haramurti/IID-mock-test/backend/web"
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
	httpapi.New(source).Register(app) // /api/* first, so it is never shadowed by the frontend
	serveFrontend(app)

	log.Printf("listening on :%s (data source: %s)", cfg.Port, cfg.DataSource)
	log.Fatal(app.Listen(":" + cfg.Port))
}

// serveFrontend serves the embedded React build. Unknown paths fall back to
// index.html so a browser refresh on any client-side URL still loads the app.
func serveFrontend(app *fiber.App) {
	if !web.Built() {
		log.Printf("frontend not embedded: run `make build` (dev: use Vite on :5173)")
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("API is running. Frontend not built: run `make build`, or `make dev` and open http://localhost:5173")
		})
		return
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(web.Dist()),
		Index:        "index.html",
		NotFoundFile: "index.html",
		MaxAge:       3600,
	}))
}

// Package httpapi wires the HTTP routes under /api.
package httpapi

import "github.com/gofiber/fiber/v2"

// Register mounts all API routes on the given router.
func Register(r fiber.Router) {
	api := r.Group("/api")
	api.Get("/health", health)
}

func health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

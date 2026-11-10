package main

// EXPECT: fiber app.Use("/api", auth) applies only to /api prefix; /admin not authed -> reported. Route registered BEFORE app.Use in fiber IS protected? (fiber: Use registers in order; earlier routes match first -> NOT protected)
import (
	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	if c.Get("Authorization") == "" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
func h(c *fiber.Ctx) error { return nil }
func main() {
	app := fiber.New()
	app.Delete("/api/early/:id", h)
	app.Use("/api", Auth)
	app.Get("/api/me", h)
	app.Post("/api/notes", h)
	app.Delete("/api/notes/:id", h)
	app.Delete("/admin/users/:id", h)
	app.Listen(":3000")
}

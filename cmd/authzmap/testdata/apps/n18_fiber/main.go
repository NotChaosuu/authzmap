package main

// EXPECT: Good enforces (fiber.ErrUnauthorized); BadAuth sets Status(401) then c.Next() -> /bad flagged only
import (
	"github.com/gofiber/fiber/v2"
)

func Good(c *fiber.Ctx) error {
	if c.Get("Authorization") == "" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
func BadAuth(c *fiber.Ctx) error {
	if c.Get("Authorization") == "" {
		c.Status(fiber.StatusUnauthorized)
	}
	return c.Next()
}
func h(c *fiber.Ctx) error { return nil }
func main() {
	app := fiber.New()
	g := app.Group("/api", Good)
	g.Delete("/users/:id", h)
	g.Post("/notes", h)
	b := app.Group("/bad", BadAuth)
	b.Delete("/users/:id", h)
	b.Post("/admin/reset", h)
	app.Listen(":3000")
}

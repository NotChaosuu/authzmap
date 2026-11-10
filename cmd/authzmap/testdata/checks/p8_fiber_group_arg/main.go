package main

// EXPECT: fiber Group("/admin", FiberGuard) -> quiet; app.Post("/ops/restart") bare -> flagged
import "github.com/gofiber/fiber/v2"

func valid(s string) bool   { return s == "k" }
func fh(c *fiber.Ctx) error { return nil }
func FiberGuard(c *fiber.Ctx) error {
	if !valid(c.Get("Authorization")) {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
func main() {
	app := fiber.New()
	admin := app.Group("/admin", FiberGuard)
	admin.Delete("/users/:id", fh)
	app.Post("/ops/restart", fh)
	app.Listen(":8080")
}

package main

// EXPECT: never rejects: fiber sets 401 status but still returns c.Next() -> auth-never-rejects MEDIUM
import "github.com/gofiber/fiber/v2"

func valid(s string) bool { return s == "x" }
func RequireAuth(c *fiber.Ctx) error {
	if !valid(c.Get("Authorization")) {
		c.Status(401)
	}
	return c.Next()
}
func h(c *fiber.Ctx) error { return nil }
func main() {
	app := fiber.New()
	g := app.Group("/api", RequireAuth)
	g.Delete("/users/:id", h)
	g.Post("/admin/reset", h)
	app.Listen(":8080")
}

package main

// EXPECT: never rejects: fiber builds fiber.ErrUnauthorized into err but returns c.Next() anyway -> auth-never-rejects MEDIUM
import "github.com/gofiber/fiber/v2"

func valid(s string) bool { return s == "x" }
func RequireAuth(c *fiber.Ctx) error {
	var err error
	if !valid(c.Get("Authorization")) {
		err = fiber.ErrUnauthorized
	}
	_ = err
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

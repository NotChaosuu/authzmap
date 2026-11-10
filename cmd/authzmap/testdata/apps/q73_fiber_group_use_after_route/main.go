package main

// EXPECT: fiber: api.Delete registered before api.Use(RequireAuth) is matched first and never reaches the auth -> DELETE /api/users/:id flagged; /api/admin/reset (after Use) quiet
import "github.com/gofiber/fiber/v2"

func valid(s string) bool { return s == "x" }
func RequireAuth(c *fiber.Ctx) error {
	if !valid(c.Get("Authorization")) {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
func h(c *fiber.Ctx) error { return nil }
func main() {
	app := fiber.New()
	api := app.Group("/api")
	api.Delete("/users/:id", h)
	api.Use(RequireAuth)
	api.Post("/admin/reset", h)
	api.Get("/me", h)
	app.Listen(":8080")
}

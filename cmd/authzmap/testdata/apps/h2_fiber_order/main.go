package main

import "github.com/gofiber/fiber/v2"

func RequireAuth(c *fiber.Ctx) error { return fiber.ErrUnauthorized }
func h(c *fiber.Ctx) error           { return nil }

func main() {
	app := fiber.New()
	app.Delete("/api/users/:id", h) // before Use: runs before the auth handler in the stack
	api := app.Group("/api")
	app.Use(RequireAuth)
	api.Post("/admin/reset", h) // after Use: auth runs first
	app.Get("/me", h)
	app.Listen(":1")
}

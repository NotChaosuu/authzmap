package main

// EXPECT: fiber middleware named RequireAuth that only logs then c.Next() -> MEDIUM; FiberGuard enforcing -> quiet
import "github.com/gofiber/fiber/v2"

func valid(s string) bool   { return s == "k" }
func fh(c *fiber.Ctx) error { return nil }
func RequireAuth(c *fiber.Ctx) error {
	println(c.Get("Authorization"))
	return c.Next()
}
func FiberGuard(c *fiber.Ctx) error {
	if !valid(c.Get("Authorization")) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}
func main() {
	app := fiber.New()
	app.Delete("/b/users/:id", RequireAuth, fh)
	app.Delete("/a/users/:id", FiberGuard, fh)
	app.Listen(":8080")
}

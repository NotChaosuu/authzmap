package main

// EXPECT: enforcing (fails closed with an empty 200 - odd, but no bypass): echo returns nil without next on a bad token -> no gaps (or unverified LOW/MEDIUM); never auth-never-rejects
import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !valid(c.Request().Header.Get("Authorization")) {
			return nil
		}
		return next(c)
	}
}
func h(c echo.Context) error { return nil }
func main() {
	e := echo.New()
	g := e.Group("/api", RequireAuth)
	g.DELETE("/users/:id", h)
	g.POST("/admin/reset", h)
	e.Start(":8080")
}

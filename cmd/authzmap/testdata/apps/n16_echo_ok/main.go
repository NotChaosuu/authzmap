package main

// EXPECT: both enforcing (return echo.ErrUnauthorized / return c.JSON(401)) -> no gaps
import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("X-Admin") != "1" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no"})
		}
		return next(c)
	}
}
func h(c echo.Context) error { return nil }
func main() {
	e := echo.New()
	g := e.Group("/api", RequireAuth)
	g.DELETE("/users/:id", h)
	a := e.Group("/admin", RequireAdmin)
	a.POST("/reset", h)
	e.Start(":8080")
}

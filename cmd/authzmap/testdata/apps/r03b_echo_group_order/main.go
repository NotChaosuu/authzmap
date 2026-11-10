package main

// EXPECT: echo group: route added before g.Use(auth) has no auth -> reported; late route authed
import (
	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func h(c echo.Context) error { return nil }
func main() {
	e := echo.New()
	g := e.Group("/api")
	g.DELETE("/early/:id", h)
	g.Use(RequireAuth)
	g.DELETE("/late/:id", h)
	g.GET("/me", h)
	g.POST("/notes", h)
	e.Start(":8080")
}

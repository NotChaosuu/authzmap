package main

// EXPECT: echo e.Use applies to ALL routes regardless of order -> no gaps. Group g.Use after a group route -> that route unprotected.
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
	e.DELETE("/users/:id", h)
	e.POST("/admin/reset", h)
	e.Use(RequireAuth)
	e.GET("/me", h)
	g := e.Group("/g")
	g.DELETE("/early/:id", h)
	g.Use(RequireAuth)
	g.DELETE("/late/:id", h)
	e.Start(":8080")
}

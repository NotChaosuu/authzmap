package main

// EXPECT: echo e.Pre(RequireAuth) runs before routing for every route -> no gaps
import (
	"github.com/labstack/echo/v4"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !valid(c.Request().Header.Get("Authorization")) {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func h(c echo.Context) error { return nil }
func main() {
	e := echo.New()
	e.Pre(RequireAuth)
	e.DELETE("/api/users/:id", h)
	e.POST("/admin/reset", h)
	e.Start(":8080")
}

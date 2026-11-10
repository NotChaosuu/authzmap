package main

// EXPECT: echo middleware only sets a header and logs, returns next(c) -> never rejects -> flagged
import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			c.Response().Header().Set("WWW-Authenticate", "Bearer")
			c.Logger().Warn("no auth ", http.StatusUnauthorized)
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

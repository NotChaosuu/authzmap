package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"strings"
)

func h(c echo.Context) error { return nil }

func main() {
	e := echo.New()
	api := e.Group("/api")
	api.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("s"),
		// skips auth for anything containing "public" -- including /api/admin/public-reset
		Skipper:                func(c echo.Context) bool { return strings.Contains(c.Path(), "admin") },
		ContinueOnIgnoredError: true,
	}))
	api.DELETE("/admin/users/:id", h)
	api.GET("/me", h)
	e.Start(":1")
}

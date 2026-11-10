package main

import "github.com/labstack/echo/v4"

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error { return echo.ErrUnauthorized }
}
func h(c echo.Context) error { return nil }

func main() {
	e := echo.New()
	g := e.Group("/api")
	g.DELETE("/users/:id", h) // before Use: echo group captures at add time
	g.POST("/admin/reset", h)
	g.Use(RequireAuth)
	g.GET("/me", h)
	g.PUT("/me", h)
	e.Start(":1")
}

package main

// EXPECT: echo middleware builds ErrUnauthorized but returns next(c) -> flagged
import (
	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		if c.Request().Header.Get("Authorization") == "" {
			err = echo.ErrUnauthorized
		}
		if err != nil {
			c.Logger().Error(err)
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

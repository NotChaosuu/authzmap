package main

// EXPECT: never rejects: echo middleware runs next(c) first, then returns ErrUnauthorized (handler already ran) -> auth-never-rejects MEDIUM
import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if !valid(c.Request().Header.Get("Authorization")) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		return err
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

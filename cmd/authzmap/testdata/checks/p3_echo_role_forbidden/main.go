package main

// EXPECT: echo middleware whose only rejection is echo.ErrForbidden / NewHTTPError(403) on a role check -> quiet; a sibling RequireAuth-named echo mw that only rate-limits with 429 -> MEDIUM
import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func valid(s string) bool { return s == "k" }
func EchoAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !valid(c.Request().Header.Get("Authorization")) {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func eh(c echo.Context) error { return nil }

type limiter struct{}

func (limiter) Allow() bool { return true }

var lim limiter

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, _ := c.Get("role").(string)
		if role != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, "admins only")
		}
		return next(c)
	}
}
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !lim.Allow() {
			return echo.NewHTTPError(http.StatusTooManyRequests)
		}
		return next(c)
	}
}
func main() {
	e := echo.New()
	e.DELETE("/admin/users/:id", eh, EchoAuth, AdminOnly)
	e.POST("/admin/reset", eh, AdminOnly)
	e.DELETE("/ops/cache", eh, RequireAuth)
	e.Start(":8080")
}

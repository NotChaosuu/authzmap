package main

// EXPECT: generated oapi routes (admin, no auth) plus hand-written authed routes -> generated routes reported, or scan INCOMPLETE. Must not be clean+complete.
import (
	"example.com/m01_oapi_echo_mixed/api"
	"github.com/labstack/echo/v4"
)

type impl struct{}

func (impl) DeleteUser(ctx echo.Context, id string) error { return nil }
func (impl) ResetAll(ctx echo.Context) error              { return nil }
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
	g := e.Group("/app", RequireAuth)
	g.GET("/me", h)
	g.POST("/notes", h)
	g.DELETE("/notes/:id", h)
	api.RegisterHandlers(e, impl{})
	e.Start(":8080")
}

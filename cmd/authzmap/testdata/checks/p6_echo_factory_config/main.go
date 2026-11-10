package main

// EXPECT: in-tree echo factory JWT(cfg) returning echo.MiddlewareFunc; cfg.Skipper always true on /b -> conditional MEDIUM; /a without skipper -> quiet
import "github.com/labstack/echo/v4"

func valid(s string) bool     { return s == "k" }
func eh(c echo.Context) error { return nil }

type JWTConfig struct{ Skipper func(echo.Context) bool }

func JWT(cfg JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Skipper != nil && cfg.Skipper(c) {
				return next(c)
			}
			if !valid(c.Request().Header.Get("Authorization")) {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
func main() {
	e := echo.New()
	a := e.Group("/a", JWT(JWTConfig{}))
	a.DELETE("/users/:id", eh)
	b := e.Group("/b", JWT(JWTConfig{Skipper: func(echo.Context) bool { return true }}))
	b.DELETE("/users/:id", eh)
	e.Start(":8080")
}

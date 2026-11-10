package main

// EXPECT: echo route with middleware spread from a slice literal holding RequireAuth -> quiet
import "github.com/labstack/echo/v4"

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
	mws := []echo.MiddlewareFunc{RequireAuth}
	e.DELETE("/api/users/:id", h, mws...)
	e.POST("/admin/reset", h, mws...)
	e.Start(":8080")
}

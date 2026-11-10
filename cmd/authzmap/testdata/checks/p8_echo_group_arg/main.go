package main

// EXPECT: echo Group with EchoAuth as arg -> quiet; sibling group without -> flagged
import "github.com/labstack/echo/v4"

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

func main() {
	e := echo.New()
	a := e.Group("/admin", EchoAuth)
	a.DELETE("/users/:id", eh)
	o := e.Group("/ops")
	o.POST("/restart", eh)
	e.Start(":8080")
}

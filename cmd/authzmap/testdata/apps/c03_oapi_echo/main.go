package main

// EXPECT: oapi-codegen echo server; generated RegisterHandlers in-tree registers via EchoRouter interface; no auth at all -> routes found and reported, or INCOMPLETE
import (
	"example.com/c03_oapi_echo/api"
	"github.com/labstack/echo/v4"
)

type impl struct{}

func (impl) DeleteUser(ctx echo.Context, id string) error { return nil }
func (impl) ResetAll(ctx echo.Context) error              { return nil }
func main() {
	e := echo.New()
	api.RegisterHandlers(e, impl{})
	e.Start(":8080")
}

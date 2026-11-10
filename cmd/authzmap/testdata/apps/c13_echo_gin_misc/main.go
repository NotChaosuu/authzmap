package main

// EXPECT: echo Add/Any/Match and gin Handle/Any/Match all found, unauthed -> reported
import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"net/http"
)

func eh(c echo.Context) error { return nil }
func gh(c *gin.Context)       {}
func main() {
	e := echo.New()
	e.Add(http.MethodDelete, "/admin/users/:id", eh)
	e.Any("/admin/any", eh)
	e.Match([]string{"POST", "PUT"}, "/admin/match", eh)
	g := gin.New()
	g.Handle("DELETE", "/gadmin/users/:id", gh)
	g.Any("/gadmin/any", gh)
	g.Match([]string{"POST"}, "/gadmin/match", gh)
	go g.Run(":8081")
	e.Start(":8080")
}

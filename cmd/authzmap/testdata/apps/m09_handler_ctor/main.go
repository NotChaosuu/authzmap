package main

// EXPECT: handler built by constructor NewAdmin(); method registered on authed group -> no gap; open one -> gap
import (
	"github.com/gin-gonic/gin"
)

type Admin struct{}

func NewAdmin() *Admin                     { return &Admin{} }
func (a *Admin) Routes(g *gin.RouterGroup) { g.DELETE("/users/:id", a.h); g.POST("/reset", a.h) }
func (a *Admin) h(c *gin.Context)          {}

type Open struct{}

func NewOpen() *Open                      { return &Open{} }
func (o *Open) Routes(g *gin.RouterGroup) { g.DELETE("/files/:id", o.h); g.POST("/reset", o.h) }
func (o *Open) h(c *gin.Context)          {}
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func main() {
	r := gin.New()
	NewAdmin().Routes(r.Group("/admin", AuthRequired()))
	NewOpen().Routes(r.Group("/open"))
	r.Run()
}

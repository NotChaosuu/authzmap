package main

// EXPECT: helper registering on groups passed from two different routers: each copy carries its caller's auth -> no gaps; /a and /b prefixes both present
import (
	"github.com/gin-gonic/gin"
)

func AuthA() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func AuthB() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Key") == "" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func h(c *gin.Context) {}
func fileRoutes(g *gin.RouterGroup) {
	up := g.Group("upload")
	{
		up.POST(":id", h)
		up.DELETE(":id", h)
	}
}
func routerA() *gin.Engine {
	r := gin.New()
	g := r.Group("/a")
	g.Use(AuthA())
	g.POST("ping", h)
	fileRoutes(g)
	return r
}
func routerB() *gin.Engine {
	r := gin.New()
	g := r.Group("/b")
	g.Use(AuthB())
	fileRoutes(g)
	return r
}
func main() { a := routerA(); b := routerB(); go a.Run(":1"); b.Run(":2") }

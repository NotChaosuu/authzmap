package main

// EXPECT: same path registered by two separate servers, each with a different enforcing auth -> no gaps
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
func h(c *gin.Context)              {}
func fileRoutes(g *gin.RouterGroup) { g.POST("upload/:id", h); g.DELETE("upload/:id", h) }
func slave() *gin.Engine {
	r := gin.New()
	g := r.Group("/api/slave")
	g.Use(AuthA())
	fileRoutes(g)
	return r
}
func master() *gin.Engine {
	r := gin.New()
	v := r.Group("/api")
	s := v.Group("slave")
	s.Use(AuthB())
	fileRoutes(s)
	return r
}
func main() { a := slave(); b := master(); go a.Run(":1"); b.Run(":2") }

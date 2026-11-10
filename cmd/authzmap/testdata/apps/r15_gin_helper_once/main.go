package main

// EXPECT: same, single caller -> no gaps
import (
	"github.com/gin-gonic/gin"
)

type Auth interface{}

func SignRequired(a Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Sign") == "" {
			c.JSON(200, gin.H{"e": 1})
			c.Abort()
			return
		}
		c.Next()
	}
}
func CacheControl() gin.HandlerFunc {
	return func(c *gin.Context) { c.Header("Cache-Control", "no-cache") }
}
func h(c *gin.Context) {}
func fileRoutes(v4 *gin.RouterGroup) {
	up := v4.Group("upload")
	{
		up.POST(":id", h)
		up.DELETE(":id", h)
	}
}
func main() {
	r := gin.New()
	v4 := r.Group("/api/v4/slave")
	v4.Use(SignRequired(nil))
	v4.Use(CacheControl())
	v4.POST("ping", h)
	fileRoutes(v4)
	r.Run()
}

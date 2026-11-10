package main

// EXPECT: per-route middleware args: DELETE has auth inline, POST passes a non-auth mw -> POST reported
import (
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func RateLimit() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func h(c *gin.Context)           {}
func main() {
	r := gin.New()
	r.DELETE("/admin/users/:id", AuthRequired(), h)
	r.POST("/admin/reset", RateLimit(), h)
	r.GET("/admin/stats", AuthRequired(), h)
	r.Run()
}

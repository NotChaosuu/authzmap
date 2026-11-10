package main

// EXPECT: v4.Use(auth) then v4.Group("upload") inline -> upload routes authed
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
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	v4 := r.Group("/api/v4")
	v4.Use(AuthRequired())
	v4.POST("ping", h)
	up := v4.Group("upload")
	up.POST(":id", h)
	up.DELETE(":id", h)
	r.Run()
}

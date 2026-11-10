package main

// EXPECT: auth on parent group passed via Group("/x", mw); child group inherits -> authed; sibling engine route unauthed -> reported
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
	api := r.Group("/api", AuthRequired())
	{
		admin := api.Group("/admin")
		admin.DELETE("/users/:id", h)
		admin.POST("/reset", h)
	}
	r.DELETE("/api/files/:id", h)
	r.Run()
}

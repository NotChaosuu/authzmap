package main

// EXPECT: c.JSON(401) without Abort; c.Next runs -> never rejects -> flagged
import (
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
		}
		c.Next()
	}
}
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	api := r.Group("/api", AuthRequired())
	api.DELETE("/users/:id", h)
	api.POST("/admin/reset", h)
	r.Run()
}

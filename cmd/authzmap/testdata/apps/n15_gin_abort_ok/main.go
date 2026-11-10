package main

// EXPECT: enforcing (AbortWithStatusJSON) -> no gaps
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
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

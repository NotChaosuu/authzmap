package main

// EXPECT: stop := c.AbortWithStatusJSON; stop(401, ...) -> quiet (enforcing)
import "github.com/gin-gonic/gin"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		stop := c.AbortWithStatusJSON
		if c.GetHeader("Authorization") == "" {
			stop(401, gin.H{"e": 1})
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

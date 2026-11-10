package main

// EXPECT: gin handlers built with append into a slice literal and spread -> quiet
import "github.com/gin-gonic/gin"

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
	chain := append([]gin.HandlerFunc{}, AuthRequired(), h)
	r.DELETE("/api/users/:id", chain...)
	r.POST("/api/admin/reset", chain...)
	r.Run()
}

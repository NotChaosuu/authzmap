package main

// EXPECT: gin middleware aborts on a missing token and never calls Next -> enforcing
import "github.com/gin-gonic/gin"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(403)
		}
	}
}
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	api := r.Group("/api")
	api.Use(AuthRequired())
	api.DELETE("/users/:id", h)
	api.POST("/admin/reset", h)
	r.Run()
}

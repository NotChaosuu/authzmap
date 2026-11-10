package main

// EXPECT: respond aborts only for 418; called with 401 the chain runs on -> never rejects MEDIUM
import "github.com/gin-gonic/gin"

func respond(c *gin.Context, code int) {
	if code == 418 {
		c.Abort()
	}
	c.JSON(code, gin.H{})
}
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			respond(c, 401)
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

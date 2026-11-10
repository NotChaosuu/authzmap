package main

// EXPECT: respond(c, 401) aborts for any code >= 400 -> quiet (enforcing)
import "github.com/gin-gonic/gin"

func respond(c *gin.Context, code int) {
	if code >= 400 {
		c.AbortWithStatus(code)
		return
	}
	c.Status(code)
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

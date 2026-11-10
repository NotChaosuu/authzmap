package main

// EXPECT: gin middleware writes 401 but neither aborts nor calls Next: gin runs the handler -> never rejects
import "github.com/gin-gonic/gin"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.String(401, "unauthorized")
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

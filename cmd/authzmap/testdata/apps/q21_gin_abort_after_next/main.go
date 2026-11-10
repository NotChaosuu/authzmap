package main

// EXPECT: never rejects: gin c.Next() then abort 401 (handler already ran) -> auth-never-rejects MEDIUM
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatus(401)
		}
	}
}
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	g := r.Group("/api", RequireAuth())
	g.DELETE("/users/:id", h)
	g.POST("/admin/reset", h)
	r.Run()
}

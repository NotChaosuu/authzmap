package main

// EXPECT: gin helper aborts on bad token; caller ignores the result but Abort already stopped the chain -> enforcing, no gaps
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func checkOrAbort(c *gin.Context) bool {
	if !valid(c.GetHeader("Authorization")) {
		c.AbortWithStatus(401)
		return false
	}
	return true
}
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkOrAbort(c)
		c.Next()
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

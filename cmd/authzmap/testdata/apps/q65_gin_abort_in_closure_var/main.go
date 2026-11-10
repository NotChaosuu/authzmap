package main

// EXPECT: enforcing: gin middleware aborts via a local closure fail := func(){ c.AbortWithStatus(401) } then return -> no gaps
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fail := func(msg string) { c.AbortWithStatusJSON(401, gin.H{"error": msg}) }
		if !valid(c.GetHeader("Authorization")) {
			fail("bad token")
			return
		}
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

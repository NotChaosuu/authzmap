package main

// EXPECT: enforcing: gin AbortWithStatusJSON(401) without return then c.Next() - Abort makes Next a no-op -> no gaps
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
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

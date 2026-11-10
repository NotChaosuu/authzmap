package main

// EXPECT: never rejects: respond() aborts only for >=500; a 401 is just c.JSON, gin continues the chain -> auth-never-rejects MEDIUM
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func respond(c *gin.Context, code int, msg string) {
	if code >= 500 {
		c.AbortWithStatusJSON(code, gin.H{"error": msg})
		return
	}
	c.JSON(code, gin.H{"error": msg})
}
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !valid(c.GetHeader("Authorization")) {
			respond(c, 401, "unauthorized")
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

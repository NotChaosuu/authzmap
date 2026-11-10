package main

// EXPECT: gin conditional where the bypass branch only returns (gin runs the chain on return) -> conditional-auth MEDIUM
import "github.com/gin-gonic/gin"

var AuthDisabled bool

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if AuthDisabled {
			return
		}
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatus(401)
			return
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

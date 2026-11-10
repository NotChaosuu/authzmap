package main

// EXPECT: gin: one route helper applied to /admin (authed) and /internal (open) groups -> /internal/* flagged; /admin quiet; no /internal/admin paths
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func h(c *gin.Context)          {}
func routes(g *gin.RouterGroup) { g.DELETE("/users/:id", h); g.POST("/reset", h) }
func main() {
	r := gin.New()
	routes(r.Group("/admin", RequireAuth()))
	routes(r.Group("/internal"))
	r.Run()
}

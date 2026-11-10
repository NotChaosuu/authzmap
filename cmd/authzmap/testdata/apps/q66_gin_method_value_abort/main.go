package main

// EXPECT: enforcing: gin abort via method value (abort := c.AbortWithStatus; abort(401)) -> no gaps; an FP never-rejects is wrong
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		abort := c.AbortWithStatus
		if !valid(c.GetHeader("Authorization")) {
			abort(401)
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

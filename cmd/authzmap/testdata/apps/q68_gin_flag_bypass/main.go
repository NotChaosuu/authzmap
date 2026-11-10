package main

// EXPECT: gin conditional: if cfg.AuthDisabled { c.Next(); return } -> conditional-auth MEDIUM
import "github.com/gin-gonic/gin"

type Config struct{ AuthDisabled bool }

var cfg Config

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.AuthDisabled {
			c.Next()
			return
		}
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatus(401)
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

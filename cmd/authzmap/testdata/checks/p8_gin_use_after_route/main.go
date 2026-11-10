package main

// EXPECT: gin: group route registered before group.Use(GinAuth()) -> that route has no auth -> flagged; route after -> quiet
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "k" }
func GinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func gh(c *gin.Context) {}

func main() {
	r := gin.New()
	g := r.Group("/admin")
	g.DELETE("/early/:id", gh)
	g.Use(GinAuth())
	g.DELETE("/late/:id", gh)
	r.Run()
}

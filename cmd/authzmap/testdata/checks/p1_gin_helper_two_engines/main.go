package main

// EXPECT: gin route helper called on an authed group of engine 1 and a bare group of engine 2 -> flagged
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

func routes(g *gin.RouterGroup) { g.DELETE("/users/:id", gh) }
func main() {
	e1 := gin.New()
	routes(e1.Group("/api", GinAuth()))
	e2 := gin.New()
	routes(e2.Group("/api"))
	go e1.Run(":8080")
	e2.Run(":8081")
}

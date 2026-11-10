package main

// EXPECT: gin auth that returns early on a nil context, then rejects anything without a valid token -> quiet (enforcing, not "credentials only when sent")
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "k" }
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c == nil {
			return
		}
		tok := c.GetHeader("Authorization")
		if tok != "" && valid(tok) {
			c.Set("user", tok)
			return
		}
		c.AbortWithStatus(401)
	}
}
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	g := r.Group("/api", AuthRequired())
	g.DELETE("/users/:id", h)
	g.POST("/admin/reset", h)
	r.Run()
}

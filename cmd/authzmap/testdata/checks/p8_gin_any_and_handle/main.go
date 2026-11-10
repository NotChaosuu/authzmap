package main

// EXPECT: gin r.Handle("DELETE", ...) and r.Any on an open group while the authed group has the same paths -> flagged
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
	v1 := r.Group("/v1", GinAuth())
	v1.DELETE("/users/:id", gh)
	v2 := r.Group("/v2")
	v2.Handle("DELETE", "/users/:id", gh)
	v2.Any("/admin/purge", gh)
	r.Run()
}

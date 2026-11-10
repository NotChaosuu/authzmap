package main

// EXPECT: gin.HandlersChain vars spread into registrations: authed chain -> quiet, a chain with only a logger -> MEDIUM
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

func Log() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func main() {
	r := gin.New()
	secure := gin.HandlersChain{Log(), GinAuth()}
	plain := gin.HandlersChain{Log()}
	r.DELETE("/a/users/:id", append(secure, gh)...)
	r.DELETE("/b/users/:id", append(plain, gh)...)
	r.Run()
}

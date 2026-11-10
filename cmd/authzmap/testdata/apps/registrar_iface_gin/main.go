package main

// EXPECT: a hand-written router interface (not oapi) over gin; routes on it carry the caller's auth
import "github.com/gin-gonic/gin"

type Router interface {
	POST(string, ...gin.HandlerFunc) gin.IRoutes
	DELETE(string, ...gin.HandlerFunc) gin.IRoutes
}

func register(r Router) { r.DELETE("/users/:id", h); r.POST("/reset", h) }
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func h(c *gin.Context) {}
func main() {
	e := gin.New()
	register(e.Group("/admin", AuthRequired()))
	register(e.Group("/open"))
	e.Run()
}

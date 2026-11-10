package main

// EXPECT: a gin middleware that ends by running another middleware's handler on the request (SignRequired(auth)(c), Cloudreve's shape) takes that verdict -> quiet
import "github.com/gin-gonic/gin"

func SignRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Sign") != secret {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func NodeRPC() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Node") == "" {
			c.JSON(200, gin.H{"code": 1})
			c.Abort()
			return
		}
		SignRequired("k")(c)
	}
}
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	g := r.Group("/api/slave", NodeRPC())
	g.DELETE("/files/:id", h)
	g.POST("/admin/reset", h)
	r.Run()
}

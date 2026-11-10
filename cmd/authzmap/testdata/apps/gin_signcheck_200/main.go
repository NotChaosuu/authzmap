package main

// EXPECT: neutral-named gin middleware aborts with 200 after a failed signature check -> treated as auth
import "github.com/gin-gonic/gin"

func verifySignature(c *gin.Context) error { return nil }
func HMACRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := verifySignature(c); err != nil {
			c.JSON(200, gin.H{"code": 40001})
			c.Abort()
			return
		}
		c.Next()
	}
}
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	slave := r.Group("/api/slave", HMACRequired())
	slave.DELETE("/files/:id", h)
	slave.POST("/admin/reset", h)
	r.GET("/api/me", h)
	r.Run()
}

package main

// EXPECT: spread chain without auth beside an authed route -> DELETE flagged MEDIUM+
import "github.com/gin-gonic/gin"

func logger() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
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
	r := gin.New()
	r.GET("/api/me", AuthRequired(), h)
	chain := []gin.HandlerFunc{logger(), h}
	r.DELETE("/api/users/:id", chain...)
	r.Run()
}

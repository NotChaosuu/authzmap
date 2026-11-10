package main

// EXPECT: r.POST registered before r.Use(auth) -> unprotected in gin -> reported; GET /me after Use authed
import (
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func h(c *gin.Context)    {}
func setup(r *gin.Engine) { r.DELETE("/users/:id", h); r.POST("/admin/reset", h) }
func main() {
	r := gin.New()
	setup(r)
	r.Use(AuthRequired())
	r.GET("/me", h)
	r.POST("/notes", h)
	r.DELETE("/notes/:id", h)
	r.Run()
}

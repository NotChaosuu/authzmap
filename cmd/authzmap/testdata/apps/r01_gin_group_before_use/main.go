package main

// EXPECT: v1 created before api.Use(auth) -> gin copies handlers at Group() -> /api/v1 routes have NO auth -> reported. /api/me after Use authed.
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
func h(c *gin.Context) {}
func main() {
	r := gin.New()
	api := r.Group("/api")
	v1 := api.Group("/v1")
	api.Use(AuthRequired())
	api.GET("/me", h)
	api.POST("/notes", h)
	api.DELETE("/notes/:id", h)
	v1.DELETE("/users/:id", h)
	v1.POST("/admin/reset", h)
	r.Run()
}

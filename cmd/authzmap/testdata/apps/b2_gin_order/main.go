package main

import "github.com/gin-gonic/gin"

func AuthRequired() gin.HandlerFunc { return func(c *gin.Context) { c.AbortWithStatus(401) } }
func h(c *gin.Context)              {}

func main() {
	r := gin.New()
	api := r.Group("/api")
	api.DELETE("/users/:id", h) // registered before Use: gin does NOT apply it
	api.POST("/admin/reset", h)
	api.Use(AuthRequired())
	api.GET("/me", h)
	api.PUT("/me", h)
	r.Run()
}

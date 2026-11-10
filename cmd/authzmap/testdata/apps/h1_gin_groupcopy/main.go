package main

import "github.com/gin-gonic/gin"

func AuthRequired() gin.HandlerFunc { return func(c *gin.Context) { c.AbortWithStatus(401) } }
func h(c *gin.Context)              {}

func main() {
	r := gin.New()
	api := r.Group("/api") // copies r's handlers now (none)
	r.Use(AuthRequired())  // too late for api
	api.DELETE("/users/:id", h)
	api.POST("/admin/reset", h)
	r.GET("/me", h)
	r.Run()
}

package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.Use(Logger())

	r.POST("/login", login)

	api := r.Group("/api", RequireAuth())
	api.GET("/me", me)
	api.GET("/orders", orders)

	// /stats is missing the per-route auth its sibling has.
	admin := r.Group("/admin")
	admin.GET("/stats", stats)
	admin.POST("/users", AdminAuth(), createUser)
}

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// an unrelated in-tree type whose Set method happens to mention 401
type apiClient struct{ status int }

func (a *apiClient) Set(k string, v any) {
	if a.status == http.StatusUnauthorized {
		a.status = 0
	}
}

// pure prefs middleware: calls gin's own ctx.Set
func prefs(c *gin.Context) {
	c.Set("lang", c.GetHeader("Accept-Language"))
	c.Next()
}

func h(c *gin.Context) {}

func main() {
	r := gin.New()
	r.Use(prefs)
	r.DELETE("/api/users/:id", h)
	r.POST("/admin/reset", h)
	r.Run()
}

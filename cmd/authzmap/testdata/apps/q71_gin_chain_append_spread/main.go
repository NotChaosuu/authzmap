package main

// EXPECT: gin route chain built with append and spread (chain...) with RequireAuth first -> no gaps; open variant without auth flagged
import "github.com/gin-gonic/gin"

func valid(s string) bool { return s == "x" }
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !valid(c.GetHeader("Authorization")) {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func h(c *gin.Context) {}
func handle(g *gin.RouterGroup, method, path string, hs ...gin.HandlerFunc) {
	chain := make([]gin.HandlerFunc, 0, len(hs)+1)
	chain = append(chain, RequireAuth())
	chain = append(chain, hs...)
	g.Handle(method, path, chain...)
}
func main() {
	r := gin.New()
	api := r.Group("/api")
	handle(api, "DELETE", "/users/:id", h)
	handle(api, "POST", "/admin/reset", h)
	api.POST("/open/admin/wipe", h)
	r.Run()
}

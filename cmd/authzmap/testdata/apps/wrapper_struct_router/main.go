package main

// EXPECT: routes go through an app's Router struct wrapping a gin group (pocket-id's shape): prefixes kept, auth from the group -> quiet; the open public group's DELETE flagged
import "github.com/gin-gonic/gin"

type Router struct{ group *gin.RouterGroup }

func NewRouter(g *gin.RouterGroup) *Router               { return &Router{group: g} }
func (r *Router) Group(p string) *Router                 { return &Router{group: r.group.Group(p)} }
func (r *Router) DELETE(p string, hs ...gin.HandlerFunc) { r.group.DELETE(p, hs...) }
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
	e := gin.New()
	api := NewRouter(e.Group("/api", AuthRequired()))
	api.Group("/users").DELETE("/:id", h)
	open := NewRouter(e.Group("/pub"))
	open.DELETE("/files/:id", h)
	e.Run()
}

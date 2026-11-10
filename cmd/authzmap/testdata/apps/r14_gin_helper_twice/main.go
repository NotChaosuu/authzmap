package main

// EXPECT: helper registers on a group passed from two routers; both callers' auth applies -> no gaps
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Auth interface{}

func SignRequired(a Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Sign") == "" {
			c.JSON(200, gin.H{"e": 1})
			c.Abort()
			return
		}
		c.Next()
	}
}
func RPCSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Sign") == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
func CacheControl() gin.HandlerFunc {
	return func(c *gin.Context) { c.Header("Cache-Control", "no-cache") }
}
func h(c *gin.Context) {}
func fileRoutes(v4 *gin.RouterGroup) {
	up := v4.Group("upload")
	{
		up.POST(":id", h)
		up.DELETE(":id", h)
	}
	f := v4.Group("file")
	{
		f.DELETE("", h)
	}
}
func slaveRouter(a Auth) *gin.Engine {
	r := gin.New()
	v4 := r.Group("/api/v4/slave")
	v4.Use(SignRequired(a))
	v4.Use(CacheControl())
	v4.POST("ping", h)
	fileRoutes(v4)
	return r
}
func masterRouter() *gin.Engine {
	r := gin.New()
	v4 := r.Group("/api/v4")
	slave := v4.Group("slave")
	slave.Use(RPCSign())
	{
		fileRoutes(slave)
	}
	return r
}
func main() {
	if len("x") == 1 {
		slaveRouter(nil).Run()
	} else {
		masterRouter().Run()
	}
}

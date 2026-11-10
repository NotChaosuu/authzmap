package main

// EXPECT: grpc-gateway behind gin.WrapH on an open Any route beside authed gin routes -> unanalysed-api MEDIUM
import (
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	gw := runtime.NewServeMux()
	r := gin.New()
	r.Group("/app", AuthRequired()).DELETE("/notes/:id", h)
	r.Any("/rpc/*path", gin.WrapH(gw))
	r.Run()
}

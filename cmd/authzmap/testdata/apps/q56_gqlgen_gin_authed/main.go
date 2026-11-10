package main

// EXPECT: gqlgen in gin under an authed group -> no gaps
import (
	"example.com/q56/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
func main() {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	r := gin.New()
	api := r.Group("/api", RequireAuth())
	api.POST("/query", gin.WrapH(srv))
	r.Run()
}

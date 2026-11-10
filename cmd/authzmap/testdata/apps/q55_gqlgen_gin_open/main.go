package main

// EXPECT: gqlgen in gin (r.POST("/query", gin.WrapH(srv))) with no auth -> unanalysed-api MEDIUM
import (
	"example.com/q55/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	r := gin.New()
	r.POST("/query", gin.WrapH(srv))
	r.Run()
}

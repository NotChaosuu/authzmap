package main

// EXPECT: gqlgen /query mounted with no auth -> reported (MEDIUM+)? mutations inside
import (
	"example.com/c07_gqlgen/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"net/http"
)

func main() {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.ListenAndServe(":8080", nil)
}

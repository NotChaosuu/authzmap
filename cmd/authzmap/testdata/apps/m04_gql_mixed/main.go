package main

// EXPECT: gqlgen /query unauthed beside authed routes -> visible at default?
import (
	"example.com/m04_gql_mixed/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/api/me", h)
		r.Post("/api/notes", h)
		r.Delete("/api/notes/{id}", h)
	})
	r.Handle("/query", srv)
	http.ListenAndServe(":8080", r)
}

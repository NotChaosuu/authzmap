package main

// EXPECT: gqlgen behind RequireAuth -> no unanalysed-api finding
import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "x", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	srv := handler.NewDefaultServer(nil)
	r := chi.NewRouter()
	r.With(RequireAuth).Handle("/graphql", srv)
	http.ListenAndServe(":8080", r)
}

package main

// EXPECT: leaf under mid; mid mounted at /x (authed) and /z (open) on root -> /z/y/* flagged, /x/y/* quiet, no stacked /x/z
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	leaf := chi.NewRouter()
	leaf.Delete("/rec/{id}", h)
	mid := chi.NewRouter()
	mid.Mount("/y", leaf)
	root := chi.NewRouter()
	root.With(RequireAuth).Mount("/x", mid)
	root.Mount("/z", mid)
	http.ListenAndServe(":8080", root)
}

package main

// EXPECT: apiRoutes is mounted only on the authed router in main; an uncalled helper (newBareRouter) also builds it open. Weakest-chain should not let dead code create a MEDIUM FP -> no gaps (or LOW)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func apiRoutes(r chi.Router)                   { r.Delete("/users/{id}", h); r.Post("/admin/reset", h) }
func newBareRouter() *chi.Mux                  { r := chi.NewRouter(); r.Route("/api", apiRoutes); return r }
func main() {
	r := chi.NewRouter()
	r.Group(func(g chi.Router) { g.Use(RequireAuth); g.Route("/api", apiRoutes) })
	http.ListenAndServe(":8080", r)
}

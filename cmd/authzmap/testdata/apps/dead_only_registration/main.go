package main

// EXPECT: the only registration of /tools/* is in an uncalled helper -> still reported (dead code is not hidden when nothing else registers it)
import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
func h(w http.ResponseWriter, r *http.Request) {}

func toolsRouter() chi.Router { r := chi.NewRouter(); r.Delete("/tools/cache", h); return r }
func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/api/me", h)
	http.ListenAndServe(":8080", r)
}

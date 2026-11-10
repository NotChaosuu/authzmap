package main

// EXPECT: routes registered through a struct of func fields (reg.post = r.Post) -> /admin/reset flagged MEDIUM or INCOMPLETE; never silent
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

type reg struct {
	post, del func(string, http.HandlerFunc)
}

func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/api/me", h)
	g := reg{post: r.Post, del: r.Delete}
	g.post("/admin/reset", h)
	g.del("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}

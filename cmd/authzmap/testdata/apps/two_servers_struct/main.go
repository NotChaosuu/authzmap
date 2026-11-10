package main

// EXPECT: same router as Handler of two http.Server values, one wrapped -> flagged
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

func main() {
	r := chi.NewRouter()
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	internal := &http.Server{Addr: ":9000", Handler: r}
	public := &http.Server{Addr: ":443", Handler: RequireAuth(r)}
	go internal.ListenAndServe()
	public.ListenAndServe()
}

package main

// EXPECT: whole mux wrapped with auth only if cfg flag -> conditional-auth MEDIUM
import (
	"net/http"
	"os"
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
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /admin/users/{id}", h)
	mux.HandleFunc("POST /admin/reset", h)
	var handler http.Handler = mux
	if os.Getenv("AUTH") == "on" {
		handler = RequireAuth(handler)
	}
	http.ListenAndServe(":8080", handler)
}

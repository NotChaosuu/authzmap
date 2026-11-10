package main

// EXPECT: &http.Server{Handler: RequireAuth(mux)} -> authed, no gaps
import (
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
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /admin/users/{id}", h)
	mux.HandleFunc("POST /admin/reset", h)
	srv := &http.Server{Addr: ":8080"}
	srv.Handler = RequireAuth(mux)
	srv.ListenAndServe()
}

package main

// EXPECT: auth wraps the mux in one branch, the bare mux is served in the other -> conditional
import (
	"net/http"
	"os"
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
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /admin/users/{id}", h)
	mux.HandleFunc("POST /admin/reset", h)
	var srv http.Handler = mux
	if os.Getenv("PUBLIC") == "" {
		srv = RequireAuth(mux)
	}
	http.ListenAndServe(":8080", srv)
}

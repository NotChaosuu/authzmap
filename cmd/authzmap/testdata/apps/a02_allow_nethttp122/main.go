package main

// EXPECT: POST/PUT/DELETE on /static/ and /assets/ with Go 1.22 patterns are NOT public
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
	mux.Handle("GET /api/me", RequireAuth(http.HandlerFunc(h)))
	mux.Handle("POST /api/notes", RequireAuth(http.HandlerFunc(h)))
	mux.Handle("DELETE /api/notes/{id}", RequireAuth(http.HandlerFunc(h)))
	mux.HandleFunc("GET /static/{path...}", h)
	mux.HandleFunc("POST /static/{file}", h)
	mux.HandleFunc("DELETE /assets/{path...}", h)
	mux.HandleFunc("PUT /assets/{name}", h)
	mux.HandleFunc("/static/", h)
	mux.HandleFunc("POST /admin/{path...}", h)
	mux.HandleFunc("example.com/admin/", h)
	mux.HandleFunc("DELETE /api/notes/{id}/{$}", h)
	http.ListenAndServe(":8080", mux)
}

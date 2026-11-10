package main

// EXPECT: ogen server (no security handler) mounted at /api/ beside authed routes -> reported at default or INCOMPLETE
import (
	"example.com/m03_ogen_mixed/oas"
	"net/http"
)

type hh struct{}

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
	srv, _ := oas.NewServer(hh{})
	mux := http.NewServeMux()
	mux.Handle("GET /app/me", RequireAuth(http.HandlerFunc(h)))
	mux.Handle("POST /app/notes", RequireAuth(http.HandlerFunc(h)))
	mux.Handle("DELETE /app/notes/{id}", RequireAuth(http.HandlerFunc(h)))
	mux.Handle("/api/", http.StripPrefix("/api", srv))
	http.ListenAndServe(":8080", mux)
}

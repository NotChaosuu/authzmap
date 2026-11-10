package main

// EXPECT: registrar interface with only HandleFunc(pattern, func): admin routes via open mux flagged MEDIUM; authed mux quiet
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

type Registrar interface {
	HandleFunc(pattern string, fn func(http.ResponseWriter, *http.Request))
}

func adminRoutes(m Registrar) {
	m.HandleFunc("DELETE /admin/users/{id}", h)
	m.HandleFunc("POST /admin/reset", h)
}
func main() {
	open := http.NewServeMux()
	adminRoutes(open)
	authed := chi.NewRouter()
	authed.Use(RequireAuth)
	adminRoutes(authed)
	go http.ListenAndServe(":9090", authed)
	http.ListenAndServe(":8080", open)
}

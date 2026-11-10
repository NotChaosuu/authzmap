package main

// EXPECT: Mount("/admin", RequireAuth(sub)) authed; Handle("/ops/*", StripPrefix(RequireAuth(sub2))) authed; With(auth).Route authed; /open/* not authed
import (
	"github.com/go-chi/chi/v5"
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
func sub() chi.Router {
	s := chi.NewRouter()
	s.Delete("/users/{id}", h)
	s.Post("/reset", h)
	return s
}
func main() {
	r := chi.NewRouter()
	r.Mount("/admin", RequireAuth(sub()))
	r.Handle("/ops/*", http.StripPrefix("/ops", RequireAuth(sub())))
	r.With(RequireAuth).Route("/billing", func(r chi.Router) { r.Delete("/cards/{id}", h); r.Post("/charge", h) })
	r.Route("/open", func(r chi.Router) { r.Delete("/admin/users/{id}", h) })
	http.ListenAndServe(":8080", r)
}

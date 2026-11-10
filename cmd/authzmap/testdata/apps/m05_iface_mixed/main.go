package main

// EXPECT: modules register via interface; admin module outside auth, notes module inside -> admin reported at default
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Module interface{ Routes(r chi.Router) }
type adminMod struct{}

func (adminMod) Routes(r chi.Router) { r.Delete("/admin/users/{id}", h); r.Post("/admin/reset", h) }

type noteMod struct{}

func (noteMod) Routes(r chi.Router) {
	r.Get("/api/notes", h)
	r.Post("/api/notes", h)
	r.Delete("/api/notes/{id}", h)
}
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
	r := chi.NewRouter()
	public := []Module{adminMod{}}
	private := []Module{noteMod{}}
	for _, m := range public {
		m.Routes(r)
	}
	r.Group(func(g chi.Router) {
		g.Use(RequireAuth)
		for _, m := range private {
			m.Routes(g)
		}
	})
	http.ListenAndServe(":8080", r)
}

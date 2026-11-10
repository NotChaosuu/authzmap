package main

// EXPECT: routes registered through an interface type; both groups visible; /admin unauthed reported
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Registrar interface {
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
}
type Module interface{ Routes(r Registrar) }
type adminMod struct{}

func (adminMod) Routes(r Registrar) { r.Delete("/admin/users/{id}", h); r.Post("/admin/reset", h) }

type noteMod struct{}

func (noteMod) Routes(r Registrar) { r.Get("/api/notes", h); r.Post("/api/notes", h) }
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
	mods := []Module{adminMod{}}
	for _, m := range mods {
		m.Routes(r)
	}
	r.Group(func(g chi.Router) { g.Use(RequireAuth); noteMod{}.Routes(g) })
	http.ListenAndServe(":8080", r)
}

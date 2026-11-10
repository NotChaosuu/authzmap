package main

// EXPECT: h := &NotesHandler{}; h.Register(g) inside auth group -> authed; admin.Register(r) outside -> gap at default
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type NotesHandler struct{ db string }

func (n *NotesHandler) Register(r chi.Router) {
	r.Get("/api/notes", n.list)
	r.Post("/api/notes", n.list)
	r.Delete("/api/notes/{id}", n.list)
}
func (n *NotesHandler) list(w http.ResponseWriter, r *http.Request) {}

type AdminHandler struct{}

func (a *AdminHandler) Register(r chi.Router) {
	r.Delete("/admin/users/{id}", a.del)
	r.Post("/admin/reset", a.del)
}
func (a *AdminHandler) del(w http.ResponseWriter, r *http.Request) {}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	notes := &NotesHandler{db: "x"}
	admin := &AdminHandler{}
	admin.Register(r)
	r.Group(func(g chi.Router) { g.Use(RequireAuth); notes.Register(g) })
	http.ListenAndServe(":8080", r)
}

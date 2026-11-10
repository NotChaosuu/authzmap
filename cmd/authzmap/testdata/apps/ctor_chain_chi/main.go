package main

// EXPECT: handlers.New(db).Mount(r.With(auth)) authed; the same on the bare router is a gap
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Notes struct{ db string }

func NewNotes(db string) *Notes                             { return &Notes{db: db} }
func (n *Notes) Mount(r chi.Router)                         { r.Delete("/notes/{id}", n.del); r.Post("/notes", n.del) }
func (n *Notes) del(w http.ResponseWriter, r *http.Request) {}

type Admin struct{}

func (Admin) Mount(r chi.Router)               { r.Post("/admin/reset", h) }
func h(w http.ResponseWriter, r *http.Request) {}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "x", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	NewNotes("db").Mount(r.With(RequireAuth))
	(&Admin{}).Mount(r)
	http.ListenAndServe(":8080", r)
}

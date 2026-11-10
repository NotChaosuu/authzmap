package main

// EXPECT: module interface with an unusual method set (Mount(r chi.Router) + Name()) ranged in a slice; open module flagged, authed quiet
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

type Module interface {
	Name() string
	Mount(r chi.Router)
}
type billing struct{}

func (billing) Name() string       { return "billing" }
func (billing) Mount(r chi.Router) { r.Delete("/billing/cards/{id}", h); r.Post("/billing/refund", h) }

type users struct{}

func (users) Name() string       { return "users" }
func (users) Mount(r chi.Router) { r.With(RequireAuth).Delete("/users/{id}", h) }
func main() {
	r := chi.NewRouter()
	mods := map[string]Module{"b": billing{}, "u": users{}}
	for _, m := range mods {
		m.Mount(r)
	}
	http.ListenAndServe(":8080", r)
}

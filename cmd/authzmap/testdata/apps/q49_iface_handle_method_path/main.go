package main

// EXPECT: in-tree interface Handle(method, path string, h http.HandlerFunc) wrapping gorilla; open admin routes MEDIUM (or INCOMPLETE), never silent
import (
	"github.com/gorilla/mux"
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

type Router interface {
	Handle(method, path string, h http.HandlerFunc)
}
type gr struct{ m *mux.Router }

func (g gr) Handle(method, path string, h http.HandlerFunc) { g.m.HandleFunc(path, h).Methods(method) }
func admin(r Router)                                        { r.Handle("DELETE", "/admin/users/{id}", h); r.Handle("POST", "/admin/reset", h) }
func main() {
	m := mux.NewRouter()
	m.Handle("/api/me", RequireAuth(http.HandlerFunc(h)))
	admin(gr{m})
	http.ListenAndServe(":8080", m)
}

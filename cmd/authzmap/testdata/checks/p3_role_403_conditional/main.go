package main

// EXPECT: role-403 gate skipped when a receiver setting DisableRBAC is set -> conditional MEDIUM (the 403 role check is the auth stop)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

type authz struct{ DisableRBAC bool }

func (a authz) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.DisableRBAC {
			next.ServeHTTP(w, r)
			return
		}
		role, _ := r.Context().Value("role").(string)
		if role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	a := authz{}
	r := chi.NewRouter()
	r.Use(a.AdminOnly)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}

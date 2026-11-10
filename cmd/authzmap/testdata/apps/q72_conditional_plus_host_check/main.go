package main

// EXPECT: auth only when a config credential is set (mailpit pattern), plus an unrelated 403 Host check -> conditional-auth MEDIUM (the Host 403 must not make it "enforcing")
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Creds struct{ U, P string }

var UICredentials *Creds

func hostOK(r *http.Request) bool { return r.Host == "localhost" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hostOK(r) {
			http.Error(w, "bad host", http.StatusForbidden)
			return
		}
		if UICredentials != nil {
			u, p, ok := r.BasicAuth()
			if !ok || u != UICredentials.U || p != UICredentials.P {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}

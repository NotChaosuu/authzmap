package main

// EXPECT: receiver verifier nil -> hand on (fail open when unconfigured) -> conditional MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}

type verifier struct{}

func (*verifier) Verify(t string) bool { return t == "k" }

type Auth struct{ v *verifier }

func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.v == nil {
			next.ServeHTTP(w, r)
			return
		}
		if !a.v.Verify(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	a := &Auth{}
	r := chi.NewRouter()
	r.Use(a.Require)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}

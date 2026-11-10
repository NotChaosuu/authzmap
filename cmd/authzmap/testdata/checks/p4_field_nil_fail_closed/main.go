package main

// EXPECT: receiver store nil -> 500 (fail closed), then a real 401 -> enforcing, quiet
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}

type store struct{}

func (*store) Valid(t string) bool { return t == "k" }

type Auth struct{ store *store }

func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.store == nil {
			http.Error(w, "misconfigured", http.StatusInternalServerError)
			return
		}
		if !a.store.Valid(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	a := &Auth{store: &store{}}
	r := chi.NewRouter()
	r.Use(a.Require)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}

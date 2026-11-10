package main

// EXPECT: admin middleware appended to a base slice and spread into the builder -> quiet; the base slice alone on /api/health is a logger only
import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hay-kot/httpkit/errchain"
)

func valid(s string) bool { return s == "k" }

type app struct{}

func (a *app) mwAuthToken(next errchain.Handler) errchain.Handler {
	return errchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}
		return next.ServeHTTP(w, r)
	})
}

// named like auth, but it only logs
func (a *app) mwRequireAuth(next errchain.Handler) errchain.Handler {
	return errchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		println(r.Header.Get("Authorization"))
		return next.ServeHTTP(w, r)
	})
}

func (a *app) mwLogger(next errchain.Handler) errchain.Handler {
	return errchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		println(r.URL.Path)
		return next.ServeHTTP(w, r)
	})
}
func del() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error { return nil }
}
func main() {
	a := &app{}
	chain := errchain.New(nil)
	base := []errchain.Middleware{a.mwLogger}
	admin := append(base, a.mwAuthToken)
	r := chi.NewRouter()
	r.Delete("/api/users/{id}", chain.ToHandlerFunc(del(), admin...))
	r.Post("/admin/reset", chain.ToHandlerFunc(del(), admin...))
	r.Get("/api/health", chain.ToHandlerFunc(del(), base...))
	http.ListenAndServe(":8080", r)
}

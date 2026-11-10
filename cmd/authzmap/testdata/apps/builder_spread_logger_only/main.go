package main

// EXPECT: chain.ToHandlerFunc(del(), mws...) where mws holds only a logger -> no auth: flagged MEDIUM
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
	mws := []errchain.Middleware{a.mwLogger}
	r := chi.NewRouter()
	r.Delete("/api/users/{id}", chain.ToHandlerFunc(del(), mws...))
	r.Post("/admin/reset", chain.ToHandlerFunc(del(), mws...))
	http.ListenAndServe(":8080", r)
}

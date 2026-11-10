package main

// EXPECT: homebox shape: chain.ToHandlerFunc(handler, userMW...) where userMW = []errchain.Middleware{a.mwAuthToken} and mwAuthToken rejects -> quiet
import (
	"github.com/go-chi/chi/v5"
	"github.com/hay-kot/httpkit/errchain"
	"net/http"
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
func del() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error { return nil }
}
func main() {
	a := &app{}
	chain := errchain.New(nil)
	userMW := []errchain.Middleware{a.mwAuthToken}
	r := chi.NewRouter()
	r.Delete("/api/v1/tags/{id}", chain.ToHandlerFunc(del(), userMW...))
	r.Put("/api/v1/tags/{id}", chain.ToHandlerFunc(del(), a.mwAuthToken))
	http.ListenAndServe(":8080", r)
}

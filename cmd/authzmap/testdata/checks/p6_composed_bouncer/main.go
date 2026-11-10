package main

// EXPECT: portainer shape: AdminAccess(h) composes h = b.mwAuthenticatedUser(h); return h, where mwAuthenticatedUser rejects with 401 -> enforcing, quiet
import (
	"github.com/gorilla/mux"
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

type Bouncer struct{}

func (b *Bouncer) mwAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (b *Bouncer) mwIsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Role") != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (b *Bouncer) TeamLeaderAccess(h http.Handler) http.Handler {
	h = b.mwIsAdmin(h)
	h = b.mwAuthenticatedUser(h)
	return h
}
func main() {
	b := &Bouncer{}
	r := mux.NewRouter()
	r.Handle("/team_memberships/{id}", b.TeamLeaderAccess(http.HandlerFunc(h))).Methods("DELETE")
	r.Handle("/team_memberships/{id}", b.TeamLeaderAccess(http.HandlerFunc(h))).Methods("PUT")
	http.ListenAndServe(":8080", r)
}

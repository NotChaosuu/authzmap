package main

// EXPECT: gorilla: same route-adding helper on an authed /admin subrouter and an open /internal subrouter -> /internal/* flagged
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
func routes(s *mux.Router) {
	s.HandleFunc("/users/{id}", h).Methods("DELETE")
	s.HandleFunc("/reset", h).Methods("POST")
}
func main() {
	r := mux.NewRouter()
	a := r.PathPrefix("/admin").Subrouter()
	a.Use(RequireAuth)
	routes(a)
	i := r.PathPrefix("/internal").Subrouter()
	routes(i)
	http.ListenAndServe(":8080", r)
}

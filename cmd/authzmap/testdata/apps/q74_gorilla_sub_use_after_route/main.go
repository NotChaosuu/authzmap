package main

// EXPECT: gorilla: subrouter Use after HandleFunc still applies (middleware runs at match time) -> no gaps
import (
	"github.com/gorilla/mux"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/users/{id}", h).Methods("DELETE")
	s.HandleFunc("/admin/reset", h).Methods("POST")
	s.Use(RequireAuth)
	http.ListenAndServe(":8080", r)
}

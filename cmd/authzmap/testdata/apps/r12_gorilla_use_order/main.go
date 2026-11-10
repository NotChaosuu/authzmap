package main

// EXPECT: gorilla Use applies at match time regardless of order -> all authed, no gaps
import (
	"github.com/gorilla/mux"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/admin/users/{id}", h).Methods("DELETE")
	r.HandleFunc("/admin/reset", h).Methods("POST")
	r.Use(RequireAuth)
	http.ListenAndServe(":8080", r)
}

package main

// EXPECT: /api/admin/x registered on ROOT router, not the authed /api subrouter -> unauthed -> reported
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
	api := r.PathPrefix("/api").Subrouter()
	api.Use(RequireAuth)
	api.HandleFunc("/me", h).Methods("GET")
	api.HandleFunc("/notes", h).Methods("POST")
	api.HandleFunc("/notes/{id}", h).Methods("DELETE")
	r.HandleFunc("/api/admin/users/{id}", h).Methods("DELETE")
	http.ListenAndServe(":8080", r)
}

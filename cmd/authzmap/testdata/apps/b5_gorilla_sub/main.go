package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler { return next }
func h(w http.ResponseWriter, r *http.Request)   {}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.Use(RequireAuth)
	api.HandleFunc("/users", h).Methods("GET")
	admin := r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/users/{id}", h).Methods("DELETE")
	admin.HandleFunc("/settings", h).Methods("POST")
	r.Methods("POST").Path("/api/users/{id}/promote").HandlerFunc(h)
	http.ListenAndServe(":8080", r)
}

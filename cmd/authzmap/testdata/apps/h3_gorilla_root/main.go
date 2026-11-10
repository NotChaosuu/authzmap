package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := mux.NewRouter()
	adm := r.PathPrefix("/admin").Subrouter()
	adm.HandleFunc("/users/{id}", h).Methods("DELETE")
	r.Use(RequireAuth)
	// separate root router served on the same port by a parent mux
	top := http.NewServeMux()
	top.Handle("/", r)
	other := mux.NewRouter()
	other.HandleFunc("/admin/drop", h).Methods("POST")
	top.Handle("/admin/drop", other)
	http.ListenAndServe(":8080", top)
}

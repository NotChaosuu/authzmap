package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	root := mux.NewRouter()
	root.StrictSlash(true)
	api := root.PathPrefix("/api/").Subrouter()
	// no auth at all on this API
	api.HandleFunc("/users/{id}", h).Methods("DELETE")
	api.HandleFunc("/admin/reset", h).Methods("POST")
	http.ListenAndServe(":8080", root)
}

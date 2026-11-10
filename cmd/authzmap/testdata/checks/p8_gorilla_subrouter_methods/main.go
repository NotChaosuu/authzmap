package main

// EXPECT: gorilla: PathPrefix(/api).Subrouter() with Use(RequireAuth) -> quiet; root.HandleFunc(/internal/flush).Methods(POST) -> flagged
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

func main() {
	root := mux.NewRouter()
	api := root.PathPrefix("/api").Subrouter()
	api.Use(RequireAuth)
	api.HandleFunc("/users/{id}", h).Methods(http.MethodDelete)
	root.HandleFunc("/internal/flush", h).Methods("POST")
	http.ListenAndServe(":8080", root)
}

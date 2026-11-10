package main

// EXPECT: gorilla: one api router behind PathPrefix /a (authed, StripPrefix) and /b (open) -> /b flagged
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
	api := mux.NewRouter()
	api.HandleFunc("/users/{id}", h).Methods("DELETE")
	root := mux.NewRouter()
	root.PathPrefix("/a/").Handler(RequireAuth(http.StripPrefix("/a", api)))
	root.PathPrefix("/b/").Handler(http.StripPrefix("/b", api))
	http.ListenAndServe(":8080", root)
}

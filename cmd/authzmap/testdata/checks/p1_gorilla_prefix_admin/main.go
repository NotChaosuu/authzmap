package main

// EXPECT: gorilla root.PathPrefix("/admin/").Handler(http.StripPrefix("/admin", api)); api has GET /static/purge and POST /users -> paths must keep /admin: GET /admin/static/purge is NOT public, POST /admin/users MEDIUM
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
	api.HandleFunc("/static/purge", h).Methods("GET")
	api.HandleFunc("/users", h).Methods("POST")
	root := mux.NewRouter()
	root.Handle("/admin/ping", RequireAuth(http.HandlerFunc(h))).Methods("POST")
	root.PathPrefix("/admin/").Handler(http.StripPrefix("/admin", api))
	http.ListenAndServe(":8080", root)
}

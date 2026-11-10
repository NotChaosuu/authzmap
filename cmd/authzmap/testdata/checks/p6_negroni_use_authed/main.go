package main

// EXPECT: control in the n26 shape: n := negroni.New(); n.Use(negroni.HandlerFunc(checkAuth)); n.UseHandler(mux) -> quiet
import (
	"github.com/urfave/negroni"
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

func checkAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !valid(r.Header.Get("Authorization")) {
		w.WriteHeader(401)
		return
	}
	next(w, r)
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /admin/users/{id}", h)
	n := negroni.New()
	n.Use(negroni.HandlerFunc(checkAuth))
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}

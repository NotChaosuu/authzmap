package main

// EXPECT: negroni.HandlerFunc(fn) with fn(w, r, next): enforcing -> quiet on the authed stack; a RequireAuth-named fn that never rejects on a second stack -> MEDIUM
import (
	"github.com/urfave/negroni"
	"net/http"
)

func valid(s string) bool                      { return s == "k" }
func h(w http.ResponseWriter, r *http.Request) {}
func checkAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !valid(r.Header.Get("Authorization")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	next(w, r)
}
func RequireAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	println(r.Header.Get("Authorization"))
	next(w, r)
}
func main() {
	a := http.NewServeMux()
	a.HandleFunc("DELETE /a/users/{id}", h)
	na := negroni.New(negroni.HandlerFunc(checkAuth))
	na.UseHandler(a)
	b := http.NewServeMux()
	b.HandleFunc("DELETE /b/users/{id}", h)
	nb := negroni.New()
	nb.Use(negroni.HandlerFunc(RequireAuth))
	nb.UseHandler(b)
	go http.ListenAndServe(":8081", na)
	http.ListenAndServe(":8080", nb)
}

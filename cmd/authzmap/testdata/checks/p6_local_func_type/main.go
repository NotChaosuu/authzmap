package main

// EXPECT: next typed by a local alias `type HandlerFn = func(http.ResponseWriter, *http.Request)` / named type `type Next http.HandlerFunc`; RequireAuth never rejects -> MEDIUM; Guard enforces -> quiet
import "net/http"

func valid(s string) bool                      { return s == "k" }
func h(w http.ResponseWriter, r *http.Request) {}

type HandlerFn = func(http.ResponseWriter, *http.Request)
type Next http.HandlerFunc

func RequireAuth(next HandlerFn) HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Header.Get("Authorization")
		next(w, r)
	}
}
func Guard(next Next) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(401)
			return
		}
		next(w, r)
	}
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /b/users/{id}", RequireAuth(h))
	mux.HandleFunc("DELETE /a/users/{id}", Guard(h))
	http.ListenAndServe(":8080", mux)
}

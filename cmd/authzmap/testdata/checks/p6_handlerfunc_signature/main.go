package main

// EXPECT: func(http.HandlerFunc) http.HandlerFunc: enforcing Auth -> quiet; RequireToken of the same signature that never rejects -> MEDIUM
import "net/http"

func valid(s string) bool                      { return s == "k" }
func h(w http.ResponseWriter, r *http.Request) {}
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(401)
			return
		}
		next(w, r)
	}
}
func RequireToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = r.Header.Get("Authorization")
		next(w, r)
	}
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /a/users/{id}", Auth(h))
	mux.HandleFunc("DELETE /b/users/{id}", RequireToken(h))
	http.ListenAndServe(":8080", mux)
}

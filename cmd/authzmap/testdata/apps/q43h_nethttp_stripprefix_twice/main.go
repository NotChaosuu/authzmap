package main

// EXPECT: net/http: one sub-mux served open at /public/ and behind RequireAuth at /admin/ -> /public/* routes flagged (DELETE /public/users/{id})
import "net/http"

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	sub := http.NewServeMux()
	sub.HandleFunc("DELETE /users/{id}", h)
	sub.HandleFunc("POST /reset", h)
	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public", sub))
	mux.Handle("/admin/", RequireAuth(http.StripPrefix("/admin", sub)))
	http.ListenAndServe(":8080", mux)
}

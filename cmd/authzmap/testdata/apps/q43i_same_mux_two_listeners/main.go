package main

// EXPECT: the same mux served open on :8080 and wrapped in RequireAuth on :9090 -> its routes flagged (the open listener)
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
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /users/{id}", h)
	mux.HandleFunc("POST /admin/reset", h)
	go http.ListenAndServe(":9090", RequireAuth(mux))
	http.ListenAndServe(":8080", mux)
}

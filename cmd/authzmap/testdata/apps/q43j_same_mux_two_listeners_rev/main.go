package main

// EXPECT: as q43i with the listeners in the other order -> routes flagged
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
	go http.ListenAndServe(":8080", mux)
	http.ListenAndServe(":9090", RequireAuth(mux))
}

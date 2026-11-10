package main

// EXPECT: Go 1.22 mux: DELETE /files/{path...} open while siblings are wrapped -> flagged
import "net/http"

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
	mux := http.NewServeMux()
	mux.Handle("GET /files/{path...}", RequireAuth(http.HandlerFunc(h)))
	mux.Handle("POST /files/{path...}", RequireAuth(http.HandlerFunc(h)))
	mux.HandleFunc("DELETE /files/{path...}", h)
	http.ListenAndServe(":8080", mux)
}

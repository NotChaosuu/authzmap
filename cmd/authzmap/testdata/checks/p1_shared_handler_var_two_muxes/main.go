package main

// EXPECT: one handler value registered on two muxes; mux A served authed, mux B served open -> flagged on B
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

func deleteUser(w http.ResponseWriter, r *http.Request) {}
func main() {
	del := http.HandlerFunc(deleteUser)
	a := http.NewServeMux()
	a.Handle("DELETE /admin/users/{id}", del)
	b := http.NewServeMux()
	b.Handle("DELETE /admin/users/{id}", del)
	go http.ListenAndServe(":9000", RequireAuth(a))
	http.ListenAndServe(":9001", b)
}

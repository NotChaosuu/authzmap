package main

// EXPECT: connect handler wrapped in RequireAuth -> no unanalysed-api finding (no gaps)
import (
	"example.com/q50/gen/adminv1/adminv1connect"
	"net/http"
)

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

type srv struct{}

func main() {
	mux := http.NewServeMux()
	path, handler := adminv1connect.NewAdminServiceHandler(srv{})
	mux.Handle(path, RequireAuth(handler))
	http.ListenAndServe(":8080", mux)
}

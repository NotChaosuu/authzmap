package main

// EXPECT: chi app with authed REST routes plus an open connect service mounted via r.Handle(adminv1connect.NewAdminServiceHandler(srv{})) -> unanalysed-api MEDIUM or INCOMPLETE; never a clean exit 0
import (
	"example.com/q52b/gen/adminv1/adminv1connect"
	"github.com/go-chi/chi/v5"
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

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Delete("/api/users/{id}", h)
	r.Handle(adminv1connect.NewAdminServiceHandler(srv{}))
	http.ListenAndServe(":8080", r)
}

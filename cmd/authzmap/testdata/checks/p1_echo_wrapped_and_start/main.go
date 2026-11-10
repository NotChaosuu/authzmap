package main

// EXPECT: echo instance served through RequireAuth on one port and via e.Start (open) on another -> flagged
import (
	"github.com/labstack/echo/v4"
	"net/http"
)

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

func eh(c echo.Context) error { return nil }
func main() {
	e := echo.New()
	e.DELETE("/admin/users/:id", eh)
	go http.ListenAndServe(":9090", RequireAuth(e))
	e.Start(":8080")
}

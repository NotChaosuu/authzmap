package main

// EXPECT: gin engine served via http.ListenAndServe(RequireAuth(e)) and also e.Run(":8080") (open) -> flagged
import (
	"github.com/gin-gonic/gin"
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

func gh(c *gin.Context) {}
func main() {
	e := gin.New()
	e.DELETE("/admin/users/:id", gh)
	go http.ListenAndServe(":9090", RequireAuth(e))
	e.Run(":8080")
}

package main

// EXPECT: fiber app served through adaptor + RequireAuth on :9 and app.Listen (open) -> flagged
import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
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

func fh(c *fiber.Ctx) error { return nil }
func main() {
	app := fiber.New()
	app.Delete("/admin/users/:id", fh)
	go http.ListenAndServe(":9", RequireAuth(adaptor.FiberApp(app)))
	app.Listen(":8080")
}

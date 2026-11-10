package main

// EXPECT: allowlist: only the genuinely public ones vanish (see per-route notes in report); DELETE /health, DELETE /webhooks/{id}, POST /admin/login, POST /static/upload, /status/../admin, /api/login/../users, POST /robots.txt, DELETE /register stay
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/api/me", h)
	r.Get("/health", h)
	r.Delete("/health", h)
	r.Post("/webhooks/github", h)
	r.Delete("/webhooks/{id}", h)
	r.Get("/api/info", h)
	r.Get("/api/system/info", h)
	r.Post("/login", h)
	r.Post("/admin/login", h)
	r.Get("/static/{path}", h)
	r.Post("/static/upload", h)
	r.Get("/metrics", h)
	r.Get("/public/{id}", h)
	r.Post("/api/public/lists", h)
	r.Get("/api/v1/auth/token", h)
	r.Post("/oauth/{provider}/callback", h)
	r.Delete("/oauth/{provider}", h)
	r.Get("/status/../admin", h)
	r.Get("/.well-known/openid-configuration", h)
	r.Post("/password/reset", h)
	r.Get("/api/login/../users", h)
	r.Get("/robots.txt", h)
	r.Post("/robots.txt", h)
	r.Get("/version/{id}/download", h)
	r.Get("/swagger/{path}", h)
	r.Post("/register", h)
	r.Delete("/register", h)
	http.ListenAndServe(":8080", r)
}

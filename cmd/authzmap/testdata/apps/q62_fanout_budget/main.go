package main

// EXPECT: fan-out helper tree (3^9 call paths per module, 50 modules) -> walk budget hit must be reported INCOMPLETE (exit 3 if no gaps), open route /open/admin/reset still flagged
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
func f0(r chi.Router)                          { r.Get("/f0", h); f1(r); f1(r); f1(r) }
func f1(r chi.Router)                          { r.Get("/f1", h); f2(r); f2(r); f2(r) }
func f2(r chi.Router)                          { r.Get("/f2", h); f3(r); f3(r); f3(r) }
func f3(r chi.Router)                          { r.Get("/f3", h); f4(r); f4(r); f4(r) }
func f4(r chi.Router)                          { r.Get("/f4", h); f5(r); f5(r); f5(r) }
func f5(r chi.Router)                          { r.Get("/f5", h); f6(r); f6(r); f6(r) }
func f6(r chi.Router)                          { r.Get("/f6", h); f7(r); f7(r); f7(r) }
func f7(r chi.Router)                          { r.Get("/f7", h); f8(r); f8(r); f8(r) }
func f8(r chi.Router)                          { r.Get("/f8", h); f9(r); f9(r); f9(r) }
func f9(r chi.Router)                          {}
func main() {
	r := chi.NewRouter()
	r.Group(func(g chi.Router) {
		g.Use(RequireAuth)
		g.Route("/m0", f0)
		g.Route("/m1", f0)
		g.Route("/m2", f0)
		g.Route("/m3", f0)
		g.Route("/m4", f0)
		g.Route("/m5", f0)
		g.Route("/m6", f0)
		g.Route("/m7", f0)
		g.Route("/m8", f0)
		g.Route("/m9", f0)
		g.Route("/m10", f0)
		g.Route("/m11", f0)
		g.Route("/m12", f0)
		g.Route("/m13", f0)
		g.Route("/m14", f0)
		g.Route("/m15", f0)
		g.Route("/m16", f0)
		g.Route("/m17", f0)
		g.Route("/m18", f0)
		g.Route("/m19", f0)
		g.Route("/m20", f0)
		g.Route("/m21", f0)
		g.Route("/m22", f0)
		g.Route("/m23", f0)
		g.Route("/m24", f0)
		g.Route("/m25", f0)
		g.Route("/m26", f0)
		g.Route("/m27", f0)
		g.Route("/m28", f0)
		g.Route("/m29", f0)
		g.Route("/m30", f0)
		g.Route("/m31", f0)
		g.Route("/m32", f0)
		g.Route("/m33", f0)
		g.Route("/m34", f0)
		g.Route("/m35", f0)
		g.Route("/m36", f0)
		g.Route("/m37", f0)
		g.Route("/m38", f0)
		g.Route("/m39", f0)
		g.Route("/m40", f0)
		g.Route("/m41", f0)
		g.Route("/m42", f0)
		g.Route("/m43", f0)
		g.Route("/m44", f0)
		g.Route("/m45", f0)
		g.Route("/m46", f0)
		g.Route("/m47", f0)
		g.Route("/m48", f0)
		g.Route("/m49", f0)
	})
	r.Post("/open/admin/reset", h)
	http.ListenAndServe(":8080", r)
}

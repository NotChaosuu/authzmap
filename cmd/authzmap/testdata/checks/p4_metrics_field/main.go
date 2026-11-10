package main

// EXPECT: receiver metrics counter incremented under `if m.metrics != nil` on the failure path -> enforcing, quiet
import (
	"github.com/go-chi/chi/v5"
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

type counter struct{}

func (*counter) Inc() {}

type MW struct {
	metrics *counter
	realm   string
}

func (m *MW) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, p, ok := r.BasicAuth(); !ok || !valid(p) {
			if m.metrics != nil {
				m.metrics.Inc()
			}
			if m.realm != "" {
				w.Header().Set("WWW-Authenticate", "Basic realm="+m.realm)
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	m := &MW{}
	r := chi.NewRouter()
	r.Use(m.Auth)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}

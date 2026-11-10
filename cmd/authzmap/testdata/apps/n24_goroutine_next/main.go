package main

// EXPECT: 401 written + return happen in a goroutine (return exits goroutine only) -> next runs -> flagged
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(t string) bool { return t == "s3cret" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		go func() {
			defer close(done)
			if !valid(r.Header.Get("Authorization")) {
				http.Error(w, "unauthorized", 401)
				return
			}
		}()
		<-done
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}

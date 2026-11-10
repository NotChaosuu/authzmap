package main

// EXPECT: middleware takes (store, inner http.Handler); enforcing -> quiet; RequireAuth(store, inner) that never rejects -> MEDIUM
import "net/http"

func h(w http.ResponseWriter, r *http.Request) {}

type Store struct{}

func (*Store) Ok(t string) bool { return t == "k" }
func Guard(s *Store, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.Ok(r.Header.Get("Authorization")) {
			w.WriteHeader(401)
			return
		}
		inner.ServeHTTP(w, r)
	})
}
func RequireAuth(s *Store, inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = s.Ok(r.Header.Get("Authorization"))
		inner.ServeHTTP(w, r)
	})
}
func main() {
	s := &Store{}
	mux := http.NewServeMux()
	mux.Handle("DELETE /a/users/{id}", Guard(s, http.HandlerFunc(h)))
	mux.Handle("DELETE /b/users/{id}", RequireAuth(s, http.HandlerFunc(h)))
	http.ListenAndServe(":8080", mux)
}

package main

// EXPECT: 13 credential stops each with an else-if credential side branch (>4096 paths) on an unknown-named gate -> unverified MEDIUM or enforcing; never silent-complete trust by name (here the name is neutral so either quiet-enforcing or unverified is acceptable; must not crash/timeout)
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

func Gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-C0") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D0") != "" {
			println(0)
		}
		if r.Header.Get("X-C1") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D1") != "" {
			println(1)
		}
		if r.Header.Get("X-C2") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D2") != "" {
			println(2)
		}
		if r.Header.Get("X-C3") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D3") != "" {
			println(3)
		}
		if r.Header.Get("X-C4") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D4") != "" {
			println(4)
		}
		if r.Header.Get("X-C5") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D5") != "" {
			println(5)
		}
		if r.Header.Get("X-C6") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D6") != "" {
			println(6)
		}
		if r.Header.Get("X-C7") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D7") != "" {
			println(7)
		}
		if r.Header.Get("X-C8") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D8") != "" {
			println(8)
		}
		if r.Header.Get("X-C9") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D9") != "" {
			println(9)
		}
		if r.Header.Get("X-C10") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D10") != "" {
			println(10)
		}
		if r.Header.Get("X-C11") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D11") != "" {
			println(11)
		}
		if r.Header.Get("X-C12") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D12") != "" {
			println(12)
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(Gate)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}

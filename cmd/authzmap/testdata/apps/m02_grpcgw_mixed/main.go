package main

// EXPECT: grpc-gateway (no auth) mounted beside authed chi routes -> gateway reported at default, or INCOMPLETE. Must not be clean+complete.
import (
	"context"
	pb "example.com/m02_grpcgw_mixed/gen/adminpb"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type srv struct{}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	gw := runtime.NewServeMux()
	pb.RegisterAdminServiceHandlerServer(context.Background(), gw, srv{})
	gw.HandlePath("DELETE", "/v1/admin/users/{id}", func(w http.ResponseWriter, r *http.Request, p map[string]string) {})
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/app/me", h)
		r.Post("/app/notes", h)
		r.Delete("/app/notes/{id}", h)
	})
	r.Mount("/v1", gw)
	http.ListenAndServe(":8080", r)
}

package main

// EXPECT: GET/POST with a runtime prefix on an unknown receiver -> listed as unseen, scan incomplete
import "net/http"

type server struct{ base string }

func (s *server) routes(mux interface {
	GET(string, http.HandlerFunc)
	POST(string, http.HandlerFunc)
}) {
	mux.POST(s.base+"/admin/reset", nil)
}
func main() {}

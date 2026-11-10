package scan

import (
	"strings"
	"testing"
)

// A router passed to a helper is resolved by type, with the caller's chain.
// The fixture builds offline through a local replace of chi.
func TestScanTypedCrossFunction(t *testing.T) {
	res, err := ScanTyped("testdata/typed_chi")
	if err != nil {
		t.Fatalf("ScanTyped: %v", err)
	}
	if !res.Complete() {
		t.Fatalf("fixture should type-check cleanly: %+v", res.Problems)
	}
	got := routeSet(res.Routes)
	chainIs(t, got, "GET /health", "Logger")
	chainIs(t, got, "GET /admin/stats", "Logger")
}

// A module that doesn't build is reported through pkg.Errors, not silently
// skipped.
func TestScanTypedReportsBuildErrors(t *testing.T) {
	dir := tree(t, map[string]string{
		"go.mod": "module example.com/broken\n\ngo 1.22\n",
		"main.go": `package main
import "net/http"
func main() {
	http.HandleFunc("/x", h)
	var n int = "not an int"
}
func h(w http.ResponseWriter, r *http.Request) {}`,
	})
	res, err := ScanTyped(dir)
	if err != nil {
		t.Fatalf("ScanTyped: %v", err)
	}
	if res.Complete() {
		t.Fatal("a type error must make the scan incomplete")
	}
	found := false
	for _, p := range res.Problems {
		if strings.Contains(p.Message, "cannot use") || strings.Contains(p.Message, "not an int") {
			found = true
		}
	}
	if !found {
		t.Errorf("problems = %+v", res.Problems)
	}
	if len(res.Routes) != 1 {
		t.Errorf("routes in the ill-typed package should still be read: %d", len(res.Routes))
	}
}

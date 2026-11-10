package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const gapApp = `package main
import "github.com/go-chi/chi/v5"
func main() {
	r := chi.NewRouter()
	r.Route("/admin", func(r chi.Router) {
		r.With(RequireAuth).Post("/users", h)
		r.Get("/stats", h)
	})
}`

func app(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, src := range files {
		p := filepath.Join(dir, name)
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func runCLI(args ...string) (int, string, string) {
	var out, errb bytes.Buffer
	code := run(args, &out, &errb)
	return code, out.String(), errb.String()
}

// Flags after the path must still apply; dropping -exit-code there would
// silently turn the CI gate off.
func TestFlagsAfterPath(t *testing.T) {
	dir := app(t, map[string]string{"main.go": gapApp})
	for _, args := range [][]string{
		{"scan", dir, "-exit-code"},
		{"scan", "-exit-code", dir},
		{"-exit-code", "scan", dir},
		{dir, "-exit-code"},
	} {
		if code, _, errs := runCLI(args...); code != exitGaps {
			t.Errorf("%v: exit %d, want %d (%s)", args, code, exitGaps, errs)
		}
	}
	_, outLow, _ := runCLI("scan", dir, "-min", "high")
	_, outHigh, _ := runCLI("scan", "-min", "high", dir)
	if outLow != outHigh {
		t.Error("-min after the path must apply too")
	}
}

func TestBadArguments(t *testing.T) {
	dir := app(t, map[string]string{"main.go": gapApp})
	for _, args := range [][]string{
		{"scan", dir, "extra"},
		{"scan", dir, "-min", "bogus"},
		{"scan", dir, "-format", "yaml"},
		{"scan", filepath.Join(dir, "missing")},
		{"scan", t.TempDir()},
	} {
		if code, _, _ := runCLI(args...); code != exitError {
			t.Errorf("%v: exit %d, want %d", args, code, exitError)
		}
	}
}

// A syntax error in the router file is loud and fails the gate.
func TestIncompleteScanFailsGate(t *testing.T) {
	clean := `package main
import "github.com/go-chi/chi/v5"
func main() { r := chi.NewRouter(); r.Use(RequireAuth); r.Get("/me", h) }
`
	dir := app(t, map[string]string{"main.go": clean + "func zzBroken( {\n"})
	code, out, errs := runCLI("scan", dir, "-exit-code")
	if code != exitIncomplete {
		t.Errorf("exit %d, want %d", code, exitIncomplete)
	}
	if !strings.Contains(errs, "INCOMPLETE") || !strings.Contains(out, "INCOMPLETE") {
		t.Errorf("stdout=%q stderr=%q", out, errs)
	}
	if code, _, _ := runCLI("scan", dir); code != exitOK {
		t.Errorf("without -exit-code the scan still exits 0, got %d", code)
	}
}

// Finding no routes must not read as "no gaps".
func TestNoRoutesIsNotClean(t *testing.T) {
	dir := app(t, map[string]string{"main.go": "package main\nfunc main() {}\n"})
	code, out, _ := runCLI("scan", dir, "-exit-code")
	if code != exitIncomplete || !strings.Contains(out, "No HTTP routes found") {
		t.Errorf("exit %d out %q", code, out)
	}
}

func TestJSONEnvelope(t *testing.T) {
	dir := app(t, map[string]string{"main.go": gapApp})
	_, out, _ := runCLI("scan", dir, "-format", "json")
	var env map[string]any
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"complete", "files", "routes", "problems", "untraced", "findings"} {
		if _, ok := env[k]; !ok {
			t.Errorf("missing %q", k)
		}
	}
	if env["problems"] == nil || env["findings"] == nil || env["untraced"] == nil {
		t.Errorf("lists must be [] not null: %s", out)
	}
	if env["complete"] != true {
		t.Error("clean scan should be complete")
	}

	_, mapOut, _ := runCLI("map", dir, "-format", "json")
	var m struct {
		Routes []struct {
			Chain []string `json:"chain"`
		} `json:"routes"`
	}
	if err := json.Unmarshal([]byte(mapOut), &m); err != nil || len(m.Routes) != 2 {
		t.Fatalf("map json: %v %s", err, mapOut)
	}
}

// map -format json says why a scan is incomplete, as scan does: an unseen
// call on an unmodelled router must be listed, not just complete=false.
func TestMapJSONListsUnseen(t *testing.T) {
	dir := app(t, map[string]string{"main.go": `package main
import ("net/http"; "github.com/julienschmidt/httprouter")
func main() {
	r := httprouter.New()
	r.GET("/admin/export", nil)
	http.ListenAndServe(":8080", r)
}`})
	_, out, _ := runCLI("map", dir, "-format", "json")
	var m struct {
		Complete bool              `json:"complete"`
		Unseen   []json.RawMessage `json:"unseen"`
	}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatal(err)
	}
	if m.Complete || len(m.Unseen) == 0 {
		t.Errorf("complete=%v unseen=%d: %s", m.Complete, len(m.Unseen), out)
	}
}

func TestAuthOverrideFlag(t *testing.T) {
	dir := app(t, map[string]string{"main.go": `package main
import "github.com/go-chi/chi/v5"
func main() {
	r := chi.NewRouter()
	r.Route("/admin", func(r chi.Router) {
		r.Use(gatekeep)
		r.Get("/stats", h)
		r.Delete("/users/{id}", h)
	})
}`})
	if code, _, _ := runCLI("scan", dir, "-exit-code"); code != exitGaps {
		t.Fatalf("unknown middleware on /admin should be a gap, exit %d", code)
	}
	if code, _, _ := runCLI("scan", dir, "-exit-code", "-auth", "gatekeep"); code != exitOK {
		t.Errorf("-auth gatekeep should clear it, exit %d", code)
	}
}

func TestVersion(t *testing.T) {
	if code, out, _ := runCLI("version"); code != exitOK || out != "authzmap "+version+"\n" {
		t.Errorf("%d %q", code, out)
	}
}

// Naming a middleware both auth and not-auth is a usage error, not a silent
// win for one side.
func TestAuthNotAuthConflict(t *testing.T) {
	code, _, errs := runCLI("scan", ".", "-auth", "Gate", "-not-auth", "gate")
	if code != 2 || !strings.Contains(errs, "both -auth and -not-auth") {
		t.Fatalf("code=%d errs=%q", code, errs)
	}
}

// Auth that never rejects, auth that rejects only sometimes and auth too
// branchy to verify each get their own status in the map table and JSON,
// with a legend line under the table.
func TestMapStatuses(t *testing.T) {
	for _, c := range []struct{ dir, label, status, legend string }{
		{"testdata/apps/a1_lognext", "NO-REJECT", "never-rejects", "NO-REJECT: named like auth, but its code never rejects"},
		{"testdata/apps/q72b_conditional_no_host_check", "COND-AUTH", "conditional", "COND-AUTH: rejects only in some cases"},
		{"testdata/checks/p2_cap_switches_authname", "UNVERIFIED", "unverified", "UNVERIFIED: named like auth, too branchy to verify"},
		{"testdata/apps/base", "NO-AUTH", "no-auth", ""},
	} {
		_, out, _ := runCLI("map", c.dir)
		if !strings.Contains(out, c.label+" ") || c.legend != "" && !strings.Contains(out, c.legend) {
			t.Errorf("%s: want %s and its legend:\n%s", c.dir, c.label, out)
		}
		if c.legend == "" && strings.Contains(out, ": ") {
			t.Errorf("%s: legend without a status that needs one:\n%s", c.dir, out)
		}
		_, js, _ := runCLI("map", c.dir, "-format", "json")
		if !strings.Contains(js, `"status": "`+c.status+`"`) {
			t.Errorf("%s: json lacks status %s", c.dir, c.status)
		}
	}
	// the method column lines up whatever the status
	_, out, _ := runCLI("map", "testdata/apps/base")
	col := -1
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		i := strings.Index(line, "/")
		j := strings.LastIndex(line[:i], "  ")
		if col == -1 {
			col = j
		} else if j != col {
			t.Errorf("misaligned row %q", line)
		}
	}
}

// scan and map JSON share one header, notes included, and lists are [] not
// null.
func TestJSONHeadersMatch(t *testing.T) {
	dir := app(t, map[string]string{"main.go": `package main
import ("net/http"; "github.com/go-pkgz/routegroup")
func main() {
	g := routegroup.New(http.NewServeMux())
	g.HandleFunc("DELETE /users/{id}", h)
}`})
	clean := app(t, map[string]string{"main.go": gapApp})
	for _, d := range []string{dir, clean} {
		_, s, _ := runCLI("scan", d, "-format", "json")
		_, m, _ := runCLI("map", d, "-format", "json")
		var sj, mj map[string]any
		if json.Unmarshal([]byte(s), &sj) != nil || json.Unmarshal([]byte(m), &mj) != nil {
			t.Fatalf("json:\n%s\n%s", s, m)
		}
		for _, k := range []string{"tool", "version", "mode", "complete", "files", "partial_routes", "problems", "untraced", "unseen", "notes"} {
			if _, ok := mj[k]; !ok {
				t.Errorf("map json lacks %q", k)
			}
			if _, ok := sj[k]; !ok {
				t.Errorf("scan json lacks %q", k)
			}
			if k != "routes" && k != "partial_routes" && k != "files" && !equalJSON(sj[k], mj[k]) {
				t.Errorf("%q differs: scan %v, map %v", k, sj[k], mj[k])
			}
		}
		for _, k := range []string{"problems", "untraced", "unseen", "notes"} {
			if sj[k] == nil || mj[k] == nil {
				t.Errorf("%q must be [] not null", k)
			}
		}
	}
	_, m, _ := runCLI("map", dir, "-format", "json")
	if !strings.Contains(m, "routegroup") {
		t.Errorf("map json lacks the unmodelled-router note:\n%s", m)
	}
}

func equalJSON(a, b any) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return bytes.Equal(x, y)
}

// Every path authzmap prints uses forward slashes, on Windows too: file
// locations, problems, unseen calls and the path given on the command line.
func TestForwardSlashes(t *testing.T) {
	dir := app(t, map[string]string{
		"internal/api/routes.go": gapApp,
		"internal/api/broken.go": "package main\nfunc broken( {",
		"internal/legacy/old.go": `package main
import "github.com/julienschmidt/httprouter"
func legacy() { r := httprouter.New(); r.GET("/legacy/export", nil) }`,
	})
	arg := dir
	if runtime.GOOS == "windows" {
		arg = strings.ReplaceAll(dir, "/", `\`)
	}
	for _, args := range [][]string{
		{"scan", arg}, {"scan", arg, "-format", "json"}, {"map", arg}, {"map", arg, "-format", "json"},
		{"scan", arg, "-v", "-min", "info"}, {"scan", filepath.Join(arg, "nope")}, {"scan", filepath.Join(arg, "internal", "api", "routes.go")},
	} {
		_, out, errs := runCLI(args...)
		if strings.Contains(out+errs, `\`) {
			t.Errorf("%v: backslash in output:\n%s\n%s", args[1:], out, errs)
		}
	}
	_, out, _ := runCLI("scan", arg, "-format", "json")
	if !strings.Contains(out, `"internal/api/routes.go"`) && !strings.Contains(out, "internal/api/routes.go:") {
		t.Errorf("want internal/api/routes.go in the output:\n%s", out)
	}
}

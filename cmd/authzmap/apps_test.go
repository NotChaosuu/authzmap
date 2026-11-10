package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// appCase pins the outcome for one app under testdata/apps, each a way a
// trust decision can fail open or flood the report.
// M: reported at medium or above. R: reported at any confidence.
// Q: present and not reported at medium or above. P: present and public.
type appCase struct {
	M, R, Q, P []string
	exit       int // -1: any nonzero
	checkExit  bool
	incomplete bool
	maxTime    time.Duration
}

var write2 = []string{"DELETE /api/users/{id}", "POST /admin/reset"}

// appTables are run as one set; a name must be unique across them.
var appTables = []map[string]appCase{middlewareApps, routingApps, coverageApps, allowlistApps}

func TestApps(t *testing.T) {
	rank := map[string]int{"high": 0, "medium": 1, "low": 2, "info": 3}
	apps := map[string]appCase{}
	for _, tbl := range appTables {
		for name, c := range tbl {
			if _, dup := apps[name]; dup {
				t.Fatalf("%s is in two tables", name)
			}
			apps[name] = c
		}
	}
	dirs, err := os.ReadDir(filepath.Join("testdata", "apps"))
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range dirs {
		if _, ok := apps[d.Name()]; !ok {
			t.Errorf("testdata/apps/%s has no expectation", d.Name())
		}
	}
	for name, want := range apps {
		name, want := name, want
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join("testdata", "apps", name)
			start := time.Now()
			_, out, errs := runCLI("scan", dir, "-format", "json", "-min", "info")
			if want.maxTime > 0 && time.Since(start) > want.maxTime {
				t.Errorf("took %s", time.Since(start))
			}
			var env struct {
				Complete bool `json:"complete"`
				Findings []struct {
					Confidence string `json:"confidence"`
					Route      *struct {
						Method, Path string
					} `json:"route"`
					Routes []string `json:"routes"`
				} `json:"findings"`
			}
			if err := json.Unmarshal([]byte(out), &env); err != nil {
				t.Fatalf("json: %v\n%s\n%s", err, out, errs)
			}
			best := map[string]string{}
			for _, f := range env.Findings {
				var keys []string
				if f.Route != nil && f.Route.Method != "" {
					keys = append(keys, f.Route.Method+" "+f.Route.Path)
				}
				for _, r := range f.Routes {
					keys = append(keys, strings.SplitN(r, " (", 2)[0])
				}
				for _, k := range keys {
					if b, ok := best[k]; !ok || rank[f.Confidence] < rank[b] {
						best[k] = f.Confidence
					}
				}
			}
			_, mapOut, _ := runCLI("map", dir, "-format", "json")
			var m struct {
				Routes []struct {
					Method, Path string
					Public       bool `json:"public"`
				} `json:"routes"`
			}
			if err := json.Unmarshal([]byte(mapOut), &m); err != nil {
				t.Fatal(err)
			}
			routes := map[string]bool{}
			public := map[string]bool{}
			for _, r := range m.Routes {
				routes[r.Method+" "+r.Path] = true
				public[r.Method+" "+r.Path] = r.Public
			}
			for _, k := range want.M {
				if c := best[k]; c != "high" && c != "medium" {
					t.Errorf("%s: want medium or above, got %q", k, c)
				}
			}
			for _, k := range want.R {
				if _, ok := best[k]; !ok {
					t.Errorf("%s: want it reported", k)
				}
			}
			for _, k := range want.Q {
				if !routes[k] {
					t.Errorf("%s: route missing", k)
				}
				if c := best[k]; c == "high" || c == "medium" {
					t.Errorf("%s: reported %s, want quiet", k, c)
				}
			}
			for _, k := range want.P {
				if !public[k] {
					t.Errorf("%s: want public", k)
				}
			}
			if want.incomplete && env.Complete {
				t.Error("want complete:false")
			}
			if want.checkExit {
				code, _, _ := runCLI("scan", dir, "-exit-code")
				if (want.exit == -1 && code == 0) || (want.exit >= 0 && code != want.exit) {
					t.Errorf("exit %d, want %d", code, want.exit)
				}
			}
		})
	}
}

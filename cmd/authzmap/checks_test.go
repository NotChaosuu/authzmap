package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// appCheck is one entry of testdata/checks/checks.json, the rules each app
// there is judged by: flag (reported at medium or above), flag_any (reported
// at any level), quiet (never medium or above), incomplete, exitcode.
type appCheck struct {
	Expect     string   `json:"expect"`
	Flag       []string `json:"flag"`
	FlagAny    []string `json:"flag_any"`
	Quiet      []string `json:"quiet"`
	Incomplete *bool    `json:"incomplete"`
	ExitCode   *int     `json:"exitcode"`
}

func TestCheckedApps(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "checks", "checks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var checks map[string]appCheck
	if err := json.Unmarshal(raw, &checks); err != nil {
		t.Fatal(err)
	}
	for name, c := range checks {
		name, c := name, c
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join("testdata", "checks", name)
			type finding struct {
				Confidence string   `json:"confidence"`
				Routes     []string `json:"routes"`
				Route      *struct {
					Method, Path string
				} `json:"route"`
			}
			var env struct {
				Complete bool      `json:"complete"`
				Findings []finding `json:"findings"`
			}
			code, out, errs := runCLI("scan", dir, "-format", "json", "-exit-code")
			if err := json.Unmarshal([]byte(out), &env); err != nil {
				t.Fatalf("json: %v\n%s", err, errs)
			}
			_, allOut, _ := runCLI("scan", dir, "-format", "json", "-min", "info")
			var all struct {
				Findings []finding `json:"findings"`
			}
			_ = json.Unmarshal([]byte(allOut), &all)
			keys := func(fs []finding, strong bool) []string {
				var ks []string
				for _, f := range fs {
					if strong && f.Confidence != "medium" && f.Confidence != "high" {
						continue
					}
					if f.Route != nil && f.Route.Method != "" {
						ks = append(ks, f.Route.Method+" "+f.Route.Path)
					}
					ks = append(ks, f.Routes...)
				}
				return ks
			}
			has := func(ks []string, want string) bool {
				for _, k := range ks {
					if strings.Contains(k, want) {
						return true
					}
				}
				return false
			}
			strong := keys(env.Findings, true)
			for _, want := range c.Flag {
				if !has(strong, want) {
					t.Errorf("want %q reported at medium or above (%s)", want, c.Expect)
				}
			}
			for _, want := range c.FlagAny {
				if !has(keys(all.Findings, false), want) {
					t.Errorf("want %q reported (%s)", want, c.Expect)
				}
			}
			for _, q := range c.Quiet {
				if has(strong, q) {
					t.Errorf("%q reported at medium or above, want quiet (%s)", q, c.Expect)
				}
			}
			if c.Incomplete != nil && env.Complete == *c.Incomplete {
				t.Errorf("complete=%v (%s)", env.Complete, c.Expect)
			}
			if c.ExitCode != nil && code != *c.ExitCode {
				t.Errorf("exit %d, want %d (%s)", code, *c.ExitCode, c.Expect)
			}
		})
	}
}

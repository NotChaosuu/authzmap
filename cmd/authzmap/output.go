package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/NotChaosuu/authzmap/pkg/model"
	"github.com/NotChaosuu/authzmap/pkg/scan"
)

// filterByConfidence keeps findings at or above min.
func filterByConfidence(fs []model.Finding, min model.Confidence) []model.Finding {
	out := []model.Finding{}
	for _, f := range fs {
		if f.Confidence.Rank() <= min.Rank() {
			out = append(out, f)
		}
	}
	return out
}

type scanJSON struct {
	Tool          string          `json:"tool"`
	Version       string          `json:"version"`
	Mode          string          `json:"mode"`
	Complete      bool            `json:"complete"`
	Files         int             `json:"files"`
	Routes        int             `json:"routes"`
	PartialRoutes int             `json:"partial_routes"`
	Problems      []model.Problem `json:"problems"`
	Untraced      []model.Site    `json:"untraced"`
	Unseen        []model.Site    `json:"unseen"`
	Notes         []string        `json:"notes"`
	Findings      []model.Finding `json:"findings"`
}

// mapJSON shares scanJSON's header; routes holds the routes, not a count.
type mapJSON struct {
	Tool          string          `json:"tool"`
	Version       string          `json:"version"`
	Mode          string          `json:"mode"`
	Complete      bool            `json:"complete"`
	Files         int             `json:"files"`
	PartialRoutes int             `json:"partial_routes"`
	Problems      []model.Problem `json:"problems"`
	Untraced      []model.Site    `json:"untraced"`
	Unseen        []model.Site    `json:"unseen"`
	Notes         []string        `json:"notes"`
	Routes        []model.Route   `json:"routes"`
}

// lists returns res's lists with nil as empty, so the JSON has [] not null.
func lists(res *scan.Result) ([]model.Problem, []model.Site, []model.Site, []string) {
	p, u, s, n := res.Problems, res.Untraced, res.Unseen, res.Notes
	if p == nil {
		p = []model.Problem{}
	}
	if u == nil {
		u = []model.Site{}
	}
	if s == nil {
		s = []model.Site{}
	}
	if n == nil {
		n = []string{}
	}
	return p, u, s, n
}

func partialCount(rs []model.Route) int {
	n := 0
	for _, r := range rs {
		if r.Partial != "" {
			n++
		}
	}
	return n
}

func printFindings(w io.Writer, res *scan.Result, fs []model.Finding, hidden int, format string) {
	if format == "json" {
		problems, untraced, unseen, notes := lists(res)
		writeJSON(w, scanJSON{
			Tool: "authzmap", Version: version, Mode: res.Mode, Complete: res.Complete() && len(res.Routes) > 0,
			Files: res.Files, Routes: len(res.Routes), PartialRoutes: partialCount(res.Routes),
			Problems: problems, Untraced: untraced, Unseen: unseen, Notes: notes, Findings: fs,
		})
		return
	}
	for _, f := range fs {
		if f.Count > 0 {
			fmt.Fprintf(w, "%-6s %s (%d routes)\n", strings.ToUpper(string(f.Confidence)), f.Rule, f.Count)
			for _, r := range capList(f.Routes, 5) {
				fmt.Fprintf(w, "       %s\n", r)
			}
		} else {
			r := f.Route
			fmt.Fprintf(w, "%-6s %-7s %s %s\n", strings.ToUpper(string(f.Confidence)), r.Method, r.Path, loc(r))
		}
		fmt.Fprintf(w, "       %s\n", f.Reason)
		if len(f.SiblingsWithAuth) > 0 {
			fmt.Fprintf(w, "       siblings with auth: %s\n", strings.Join(capList(f.SiblingsWithAuth, 3), ", "))
		}
		if f.Hint != "" {
			fmt.Fprintf(w, "       hint: %s\n", f.Hint)
		}
		fmt.Fprintln(w)
	}
	s := summary(res, gapRoutes(fs))
	if hidden > 0 {
		s += fmt.Sprintf(" %d lower-confidence finding(s) hidden (use -min low or -min info).", hidden)
	}
	fmt.Fprintln(w, s)
}

// gapRoutes counts routes with a finding: a grouped or app-wide finding
// counts its routes once, not again beside their own findings.
func gapRoutes(fs []model.Finding) int {
	seen := map[string]bool{}
	for _, f := range fs {
		if f.Route.Method != "" {
			seen[f.Route.Method+" "+f.Route.Path+" ("+f.Route.File+":"+strconv.Itoa(f.Route.Line)+")"] = true
		}
		for _, r := range f.Routes {
			seen[r] = true
		}
	}
	return len(seen)
}

func summary(res *scan.Result, gaps int) string {
	where := fmt.Sprintf("%d Go file(s)", res.Files)
	switch {
	case len(res.Routes) == 0 && !res.Complete():
		return fmt.Sprintf("No HTTP routes found in %s, and the scan is INCOMPLETE: nothing was checked.", where)
	case len(res.Routes) == 0:
		return fmt.Sprintf("No HTTP routes found in %s: nothing was checked (supported: chi, gin, echo, fiber, gorilla/mux, hertz, net/http).", where)
	}
	routes := fmt.Sprintf("%d route(s)", len(res.Routes))
	if p := partialCount(res.Routes); p > 0 {
		routes += fmt.Sprintf(", %d partial", p)
	}
	var s string
	if gaps == 0 {
		s = fmt.Sprintf("No authorization gaps found across %s in %s.", routes, where)
	} else {
		s = fmt.Sprintf("%d gap(s) across %s in %s.", gaps, routes, where)
	}
	if !res.Complete() {
		s += fmt.Sprintf(" INCOMPLETE: %d problem(s), see stderr.", len(res.Problems)+len(res.Unseen))
	}
	return s
}

func capList(xs []string, n int) []string {
	if len(xs) <= n {
		return xs
	}
	return append(append([]string(nil), xs[:n]...), fmt.Sprintf("and %d more", len(xs)-n))
}

func printMap(w io.Writer, res *scan.Result, routes []model.Route, format string) {
	if format == "json" {
		if routes == nil {
			routes = []model.Route{}
		}
		problems, untraced, unseen, notes := lists(res)
		writeJSON(w, mapJSON{
			Tool: "authzmap", Version: version, Mode: res.Mode, Complete: res.Complete() && len(res.Routes) > 0,
			Files: res.Files, PartialRoutes: partialCount(routes),
			Problems: problems, Untraced: untraced, Unseen: unseen, Notes: notes, Routes: routes,
		})
		return
	}
	sort.SliceStable(routes, func(i, j int) bool {
		if routes[i].Path != routes[j].Path {
			return routes[i].Path < routes[j].Path
		}
		return routes[i].Method < routes[j].Method
	})
	// pad paths to a common column; an outlier past 60 runs over instead of
	// pushing every row right
	width := 28
	for _, r := range routes {
		if l := len(r.Path); l > width && l <= 60 {
			width = l
		}
	}
	for _, r := range routes {
		chain := strings.Join(r.Chain, " -> ")
		if chain == "" {
			chain = "-"
		}
		mark := " "
		if r.Partial != "" {
			mark = "~"
		}
		fmt.Fprintf(w, "%-10s%s %-7s %-*s (%s)\n", label[r.Status], mark, r.Method, width, r.Path, chain)
	}
	if len(routes) == 0 {
		fmt.Fprintln(w, summary(res, 0))
		return
	}
	var legend []string
	for _, st := range []string{model.StatusNeverRejects, model.StatusConditional, model.StatusUnverified} {
		for _, r := range routes {
			if r.Status == st {
				legend = append(legend, label[st]+": "+meaning[st])
				break
			}
		}
	}
	if partialCount(routes) > 0 {
		legend = append(legend, "~ partial: middleware applied where this router came from wasn't traced")
	}
	if len(legend) > 0 {
		fmt.Fprintln(w, "\n"+strings.Join(legend, "\n"))
	}
}

// label is the map table's status column. Upper case is exposed, or may be.
var label = map[string]string{
	model.StatusAuth:         "auth",
	model.StatusPublic:       "public",
	model.StatusWeak:         "weak",
	model.StatusUnverified:   "UNVERIFIED",
	model.StatusConditional:  "COND-AUTH",
	model.StatusNeverRejects: "NO-REJECT",
	model.StatusNoAuth:       "NO-AUTH",
}

var meaning = map[string]string{
	model.StatusNeverRejects: "named like auth, but its code never rejects",
	model.StatusConditional:  "rejects only in some cases (a setting, a method)",
	model.StatusUnverified:   "named like auth, too branchy to verify",
}

func loc(r model.Route) string {
	if r.File == "" {
		return ""
	}
	return fmt.Sprintf("(%s:%d)", r.File, r.Line)
}

func writeJSON(w io.Writer, v any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

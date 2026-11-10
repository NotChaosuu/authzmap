package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/NotChaosuu/authzmap/pkg/analyze"
	"github.com/NotChaosuu/authzmap/pkg/model"
	"github.com/NotChaosuu/authzmap/pkg/scan"
)

// version is set at release time with -ldflags "-X main.version=..."; a
// plain go build or go install reports the version below.
var version = "1.0.0"

// Exit codes.
const (
	exitOK         = 0
	exitGaps       = 1 // -exit-code: gaps at or above -min
	exitError      = 2 // bad usage, unreadable path, no Go files
	exitIncomplete = 3 // -exit-code: some files didn't parse/type-check, or no routes were found
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

type config struct {
	cmd      string
	path     string
	format   string
	min      model.Confidence
	opts     analyze.Options
	exitCode bool
	typed    bool
	verbose  bool
}

func run(args []string, stdout, stderr io.Writer) int {
	cfg, err := parseArgs(args, stderr)
	if errors.Is(err, flag.ErrHelp) {
		return exitOK
	}
	if err != nil {
		fmt.Fprintln(stderr, "authzmap:", err)
		return exitError
	}
	if cfg.cmd == "version" {
		fmt.Fprintln(stdout, "authzmap", version)
		return exitOK
	}

	scanFn := scan.Scan
	if cfg.typed {
		scanFn = scan.ScanTyped
	}
	res, err := scanFn(cfg.path)
	if err != nil {
		fmt.Fprintln(stderr, "authzmap:", err)
		return exitError
	}
	warn(res, stderr, cfg.verbose)

	incomplete := !res.Complete() || len(res.Routes) == 0
	switch cfg.cmd {
	case "map":
		printMap(stdout, res, analyze.Classify(res.Routes, cfg.opts), cfg.format)
		if cfg.exitCode && incomplete {
			return exitIncomplete
		}
		return exitOK
	default:
		all := analyze.Findings(res.Routes, cfg.opts)
		findings := filterByConfidence(all, cfg.min)
		warnOverrides(res, cfg, stderr)
		printFindings(stdout, res, findings, len(all)-len(findings), cfg.format)
		switch {
		case !cfg.exitCode:
			return exitOK
		case len(findings) > 0:
			return exitGaps
		case incomplete:
			return exitIncomplete
		}
		return exitOK
	}
}

// parseArgs accepts flags before, between and after the positionals:
// `authzmap scan ./svc -exit-code` and `authzmap -exit-code scan ./svc`
// mean the same thing. Go's flag package stops at the first positional, so
// it's re-run on whatever follows each one.
func parseArgs(args []string, stderr io.Writer) (*config, error) {
	cfg := &config{cmd: "scan", path: "."}
	fs := flag.NewFlagSet("authzmap", flag.ContinueOnError)
	fs.SetOutput(stderr)
	format := fs.String("format", "table", "output format: table or json")
	min := fs.String("min", "medium", "lowest confidence to report: high, medium, low or info")
	allow := fs.String("allow", "", "comma-separated public path globs, merged with the defaults")
	noDefault := fs.Bool("no-default-allow", false, "drop the built-in public allowlist")
	auth := fs.String("auth", "", "comma-separated middleware names to treat as auth checks")
	notAuth := fs.String("not-auth", "", "comma-separated middleware names that are not auth checks")
	exitCode := fs.Bool("exit-code", false, "exit 1 on gaps at or above -min, 3 on an incomplete scan (for CI)")
	typed := fs.Bool("types", false, "type-check the target (go/types); needs its dependencies to resolve")
	verbose := fs.Bool("v", false, "list every unreadable file and untraced registration")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "authzmap: find routes reachable without an authorization check, in Go source.")
		fmt.Fprintln(stderr, "\nUsage:\n  authzmap [scan|map] [path] [flags]\n  authzmap version")
		fmt.Fprintln(stderr, "\nExit codes: 0 ok, 2 usage or path error; with -exit-code also 1 gaps, 3 incomplete scan.")
		fmt.Fprintln(stderr, "\nFlags:")
		fs.PrintDefaults()
	}

	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		args = fs.Args()
		if len(args) == 0 {
			break
		}
		pos = append(pos, args[0])
		args = args[1:]
	}
	if len(pos) > 0 {
		switch pos[0] {
		case "scan", "map", "version":
			cfg.cmd = pos[0]
			pos = pos[1:]
		}
	}
	switch len(pos) {
	case 0:
	case 1:
		cfg.path = strings.TrimSuffix(strings.TrimSuffix(pos[0], "/..."), `\...`)
		// forward slashes everywhere a path is printed; Windows takes them too
		cfg.path = filepath.ToSlash(cfg.path)
		if cfg.path == "" || cfg.path == "..." {
			cfg.path = "."
		}
	default:
		return nil, fmt.Errorf("one path at a time; unexpected %q", pos[1])
	}

	switch *format {
	case "table", "json":
		cfg.format = *format
	default:
		return nil, fmt.Errorf("-format must be table or json, not %q", *format)
	}
	c, ok := model.ParseConfidence(*min)
	if !ok {
		return nil, fmt.Errorf("-min must be high, medium, low or info, not %q", *min)
	}
	cfg.min = c
	cfg.opts = analyze.Options{
		AllowGlobs: splitList(*allow), NoDefaultAllow: *noDefault,
		Auth: splitList(*auth), NotAuth: splitList(*notAuth),
	}
	for _, a := range cfg.opts.Auth {
		for _, n := range cfg.opts.NotAuth {
			if strings.EqualFold(a, n) {
				return nil, fmt.Errorf("%q is in both -auth and -not-auth", a)
			}
		}
	}
	cfg.exitCode, cfg.typed, cfg.verbose = *exitCode, *typed, *verbose
	return cfg, nil
}

// warnOverrides: an unqualified -auth name matches the last element of any
// reference, so -auth Handler also blesses cors.Handler. Say which
// qualified references it matched.
func warnOverrides(res *scan.Result, cfg *config, w io.Writer) {
	for _, name := range cfg.opts.Auth {
		if strings.Contains(name, ".") {
			continue
		}
		seen := map[string]bool{}
		var hits []string
		for _, r := range res.Routes {
			for _, m := range r.Chain {
				ref := m
				if i := strings.IndexByte(ref, '='); i >= 0 {
					ref = ref[:i]
				}
				if strings.Contains(ref, ".") && strings.EqualFold(ref[strings.LastIndexByte(ref, '.')+1:], name) && !seen[ref] {
					seen[ref] = true
					hits = append(hits, ref)
				}
			}
		}
		if len(hits) > 0 {
			fmt.Fprintf(w, "authzmap: -auth %s also matches %s; pass the qualified name to be exact\n", name, strings.Join(hits, ", "))
		}
	}
}

func splitList(s string) []string {
	var out []string
	for _, x := range strings.Split(s, ",") {
		if x = strings.TrimSpace(x); x != "" {
			out = append(out, x)
		}
	}
	return out
}

// warn puts what the scan couldn't read on stderr, so it shows even when
// stdout is piped into something else.
func warn(res *scan.Result, w io.Writer, verbose bool) {
	if n := len(res.Problems); n > 0 {
		what := "problem(s) (unparsable files, call chains cut short)"
		if res.Mode == "types" {
			what = "problem(s) loading or type-checking packages"
		}
		fmt.Fprintf(w, "authzmap: INCOMPLETE SCAN: %d %s; routes behind them may be missing.\n", n, what)
		limit := 5
		if verbose {
			limit = n
		}
		for i, p := range res.Problems {
			if i == limit {
				fmt.Fprintf(w, "  ... and %d more (-v lists them all)\n", n-limit)
				break
			}
			loc := p.File
			if loc == "" {
				loc = "(package)"
			}
			fmt.Fprintf(w, "  %s: %s\n", loc, p.Message)
		}
	}
	if n := len(res.Unseen); n > 0 {
		fmt.Fprintf(w, "authzmap: INCOMPLETE SCAN: %d route-like call(s) belong to a router authzmap doesn't model; they weren't checked.\n", n)
		limit := 5
		if verbose {
			limit = n
		}
		for i, s := range res.Unseen {
			if i == limit {
				fmt.Fprintf(w, "  ... and %d more (-v lists them all)\n", n-limit)
				break
			}
			fmt.Fprintf(w, "  %s:%d %s\n", s.File, s.Line, s.Call)
		}
	}
	if len(res.Untraced) > 0 {
		fmt.Fprintf(w, "authzmap: %d registration(s) on values not traced to a router; reported as partial.\n", len(res.Untraced))
		if verbose {
			for _, s := range res.Untraced {
				fmt.Fprintf(w, "  %s:%d %s\n", s.File, s.Line, s.Call)
			}
		}
	}
	for _, n := range res.Notes {
		fmt.Fprintln(w, "authzmap: note:", n)
	}
}

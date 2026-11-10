// Package model holds the types shared by the scanner, the analyzer and the CLI.
package model

import (
	"encoding/json"
	"strings"
)

// Confidence is how sure a finding is that the route is really exposed.
type Confidence string

const (
	High   Confidence = "high"
	Medium Confidence = "medium"
	Low    Confidence = "low"
	Info   Confidence = "info"
)

var confidenceRank = map[Confidence]int{High: 0, Medium: 1, Low: 2, Info: 3}

// Rank orders confidences, most confident first. Unknown values sort last.
func (c Confidence) Rank() int {
	if r, ok := confidenceRank[c]; ok {
		return r
	}
	return len(confidenceRank)
}

// ParseConfidence accepts high, medium, low or info (any case).
func ParseConfidence(s string) (Confidence, bool) {
	switch c := Confidence(strings.ToLower(s)); c {
	case High, Medium, Low, Info:
		return c, true
	}
	return "", false
}

// Scope kinds. Sibling inconsistency is only trusted inside an explicit group;
// a root router legitimately mixes public and authenticated routes.
const (
	ScopeRoot   = "root"   // an engine/router created with its constructor
	ScopeGroup  = "group"  // Group/Route/subrouter/mounted router
	ScopePrefix = "prefix" // net/http: grouped by first path segment
	ScopeParam  = "param"  // a router that arrived as a parameter we couldn't trace
)

// Route is one registered handler with everything known about the
// middleware in front of it.
type Route struct {
	Method    string   `json:"method"`
	Path      string   `json:"path"`
	Handler   string   `json:"handler"`
	Framework string   `json:"framework"`
	File      string   `json:"file"`
	Line      int      `json:"line"`
	Chain     []string `json:"chain"`
	Scope     string   `json:"scope"`
	ScopeKind string   `json:"scope_kind"`
	// Partial is set when middleware applied outside what was traced (by a
	// caller, or on a router whose origin wasn't found) may be missing.
	Partial string `json:"partial,omitempty"`
	// Dynamic is set when the path isn't a constant; Path then holds the
	// expression in angle brackets.
	Dynamic bool `json:"dynamic,omitempty"`
	// Enforcing lists chain entries whose own code was read and rejects
	// unauthenticated requests (a 401 or a redirect to login).
	Enforcing []string `json:"enforcing,omitempty"`
	// HandlerChecks is set when the handler's own body rejects anonymous
	// callers.
	HandlerChecks bool `json:"handler_checks_auth,omitempty"`
	// Passing lists chain entries whose own code always hands the request
	// on: they never reject anything, whatever their name.
	Passing []string `json:"passing,omitempty"`
	// Conditional lists chain entries added inside an if/loop/switch.
	Conditional []string `json:"conditional,omitempty"`
	// API names an API server mounted as one handler (grpc-gateway,
	// gqlgen, connect-rpc, a generated OpenAPI server, routes built by
	// reflection): its operations aren't analysed one by one.
	API string `json:"api,omitempty"`
	// Variants holds every chain the same registration was reached with
	// (a helper called from two routers); the route is only as protected
	// as its weakest one.
	Variants []Variant `json:"-"`
	// Unverified lists chain entries whose code was too branchy to read.
	Unverified []string `json:"unverified,omitempty"`
	// AuthElsewhere is set when another way to the same registration
	// requires auth.
	AuthElsewhere bool `json:"auth_elsewhere,omitempty"`

	// Status sums up the route in one word; see the Status constants.
	Status   string   `json:"status"`
	HasAuth  bool     `json:"has_auth"`
	WeakAuth []string `json:"weak_auth,omitempty"`
	CondAuth []string `json:"conditional_auth,omitempty"`
	// PassThrough lists middleware named like an auth check whose code
	// never rejects a request.
	PassThrough []string `json:"auth_never_rejects,omitempty"`
	Public      bool     `json:"public"`
}

// Route statuses, strongest first.
const (
	StatusAuth         = "auth"          // an auth check rejects anonymous callers
	StatusPublic       = "public"        // public by design (login, health, static)
	StatusWeak         = "weak"          // only something that may be auth (a session loader, an optional check)
	StatusUnverified   = "unverified"    // named like auth, but too branchy to verify
	StatusConditional  = "conditional"   // auth that rejects only in some cases (a setting, a method)
	StatusNeverRejects = "never-rejects" // named like auth, but its code never rejects
	StatusNoAuth       = "no-auth"       // nothing in front of it looks like auth
)

// Variant is one way a registration was reached.
type Variant struct {
	Chain, Enforcing, Passing, Conditional, Unverified []string
	HandlerChecks                                      bool
	Partial                                            string
}

// Finding is one reported gap: a route, or for app-wide and grouped rules a
// set of routes, reachable without an authorization check.
type Finding struct {
	Route            Route      `json:"route"`
	Rule             string     `json:"rule"`
	Confidence       Confidence `json:"confidence"`
	Reason           string     `json:"reason"`
	SiblingsWithAuth []string   `json:"siblings_with_auth,omitempty"`
	// Count is set on findings that summarize many routes; Routes then lists
	// them as "METHOD /path (file:line)".
	Count  int      `json:"count,omitempty"`
	Routes []string `json:"routes,omitempty"`
	Hint   string   `json:"hint,omitempty"`
}

// MarshalJSON leaves the route out of app-level findings that aren't about
// one route (no-auth-anywhere), instead of emitting an empty object.
func (f Finding) MarshalJSON() ([]byte, error) {
	type plain Finding
	if f.Route.Method == "" && f.Route.Path == "" {
		return json.Marshal(struct {
			plain
			Route *Route `json:"route,omitempty"`
		}{plain: plain(f)})
	}
	return json.Marshal(plain(f))
}

// Problem is a file or package the scan could not fully read.
type Problem struct {
	File    string `json:"file,omitempty"`
	Message string `json:"message"`
}

// Site is a source location, used for route-like calls that couldn't be
// attributed to a router.
type Site struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Call string `json:"call"`
}

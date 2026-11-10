// Package analyze turns routes into authorization-gap findings. A route that
// skips the auth its group peers use (sibling inconsistency) ranks highest; a
// route in an app with no auth anywhere ranks lowest.
package analyze

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/NotChaosuu/authzmap/pkg/authsig"
	"github.com/NotChaosuu/authzmap/pkg/model"
)

// Options are the user's adjustments to what counts as public and as auth.
type Options struct {
	// AllowGlobs are extra public path patterns merged with the defaults.
	AllowGlobs []string
	// NoDefaultAllow drops the built-in public allowlist (login, health, ...).
	NoDefaultAllow bool
	// Auth and NotAuth name middleware the vocabulary gets wrong.
	Auth    []string
	NotAuth []string
}

func (o Options) classifier() *authsig.Classifier {
	return authsig.NewClassifier(o.Auth, o.NotAuth)
}

// Classify returns the routes with HasAuth, WeakAuth and Public filled in.
func Classify(routes []model.Route, opts Options) []model.Route {
	c := opts.classifier()
	out := make([]model.Route, len(routes))
	for i, r := range routes {
		if len(r.Variants) < 2 {
			out[i] = classifyOne(r, c, opts)
			continue
		}
		// reached through several chains: judged by the weakest
		best, rank, anyAuth := model.Route{}, 99, false // 99: worse than any real rank
		for _, v := range r.Variants {
			rv := r
			rv.Chain, rv.Enforcing, rv.Passing, rv.Conditional, rv.Unverified = v.Chain, v.Enforcing, v.Passing, v.Conditional, v.Unverified
			rv.HandlerChecks, rv.Partial = v.HandlerChecks, v.Partial
			rv = classifyOne(rv, c, opts)
			if k := strength(rv); k < rank {
				best, rank = rv, k
			}
			anyAuth = anyAuth || rv.HasAuth
		}
		best.AuthElsewhere = anyAuth && !best.HasAuth
		out[i] = best
	}
	return out
}

func strength(r model.Route) int {
	switch {
	case r.HasAuth:
		return 4
	case len(r.WeakAuth) > 0:
		return 3
	case len(r.CondAuth) > 0:
		return 2
	case len(r.PassThrough) > 0:
		return 1
	}
	return 0
}

// baseName drops the note a conditional verdict appends ("RequireAuth
// [rejects only some methods]") before the name is classified.
func baseName(m string) string {
	if i := strings.Index(m, " ["); i > 0 {
		return m[:i]
	}
	return m
}

func classifyOne(r model.Route, c *authsig.Classifier, opts Options) model.Route {
	r.HasAuth, r.WeakAuth, r.CondAuth, r.PassThrough = false, nil, nil, nil
	verified := map[string]bool{}
	for _, m := range r.Enforcing {
		verified[m] = true
	}
	cond := map[string]bool{}
	for _, m := range r.Conditional {
		cond[m] = true
	}
	passing := map[string]bool{}
	for _, m := range r.Passing {
		passing[m] = true
	}
	unverified := map[string]bool{}
	for _, m := range r.Unverified {
		unverified[m] = true
	}
	for _, m := range r.Chain {
		k := c.Classify(baseName(m))
		// Code that always passes the request on isn't a check, whatever
		// it's called, unless the user said it is.
		if passing[m] && !c.Overridden(baseName(m)) {
			if k == authsig.Auth || k == authsig.Weak {
				r.PassThrough = append(r.PassThrough, m)
			}
			k = authsig.NotAuth
		}
		// code too branchy to read is unverified, whatever its name says
		if unverified[m] && !c.Overridden(baseName(m)) && k == authsig.Auth {
			k = authsig.Weak
		}
		// "Handler [rejects only when a setting is on]": its code was read
		// and checks, whatever the name
		fromCode := baseName(m) != m && k != authsig.NotAuth
		if cond[m] && (k == authsig.Auth || verified[m] && k != authsig.NotAuth || fromCode) {
			r.CondAuth = append(r.CondAuth, m)
			continue
		}
		// Code that rejects anonymous requests outranks its name, unless
		// the name says it's optional or the user said it isn't auth.
		if verified[m] && k != authsig.NotAuth && !c.Optional(baseName(m)) {
			k = authsig.Auth
		}
		switch k {
		case authsig.Auth:
			r.HasAuth = true
		case authsig.Weak:
			r.WeakAuth = append(r.WeakAuth, m)
		}
	}
	if r.HasAuth {
		r.WeakAuth, r.CondAuth = nil, nil
	}
	r.Public = isPublic(r.Method, r.Path, opts)
	r.Status = status(r)
	return r
}

func status(r model.Route) string {
	switch {
	case r.HasAuth:
		return model.StatusAuth
	case r.Public:
		return model.StatusPublic
	}
	for _, w := range r.WeakAuth {
		for _, u := range r.Unverified {
			if w == u {
				return model.StatusUnverified
			}
		}
	}
	switch {
	case len(r.WeakAuth) > 0:
		return model.StatusWeak
	case len(r.CondAuth) > 0:
		return model.StatusConditional
	case len(r.PassThrough) > 0:
		return model.StatusNeverRejects
	}
	return model.StatusNoAuth
}

// Rule names.
const (
	RuleSibling     = "sibling-inconsistency"
	RuleUnverified  = "unverified-auth"
	RuleConditional = "conditional-auth"
	RuleSensitive   = "sensitive-no-auth"
	RuleNoAuth      = "no-auth"
	RuleNoAuthAtAll = "no-auth-anywhere"
	RulePassThrough = "auth-never-rejects"
	RuleAPI         = "unanalysed-api"

	// Percentage of a route's non-public siblings that must be authed before
	// its gap counts as inconsistent; see siblingEvidence and pathSiblings.
	minLocalPct = 50
	minScopePct = 75
	minPathPct  = 75
)

// Findings reports the routes reachable without an authorization check,
// most confident first. Public (allowlisted) routes are never findings.
func Findings(routes []model.Route, opts Options) []model.Finding {
	classified := Classify(routes, opts)

	authed := map[string][]string{}
	exposed := map[string]int{} // non-public routes per scope, and per scope + first segment
	appAuth := false
	for _, r := range classified {
		if r.Public {
			continue
		}
		local, path := localKey(r), pathKey(r)
		exposed[r.Scope]++
		exposed[local]++
		exposed[path]++
		// a handler that checks the caller itself counts as a guarded close
		// sibling, so the one handler that forgot stands out
		id := r.Method + " " + r.Path
		switch {
		case r.HasAuth:
			appAuth = true
			authed[r.Scope] = append(authed[r.Scope], id)
			authed[local] = append(authed[local], id)
			authed[path] = append(authed[path], id)
			if vk := versionKey(r.Path); vk != "" {
				authed[vk] = append(authed[vk], id)
			}
		case r.HandlerChecks, r.AuthElsewhere, len(r.PassThrough) > 0, len(r.CondAuth) > 0, len(r.WeakAuth) > 0:
			// a handler that checks the caller, or middleware meant as auth
			// that falls short, makes the app one that authenticates; none
			// of it is sibling evidence
			appAuth = true
		}
	}

	var findings []model.Finding
	unguarded := 0
	for _, r := range classified {
		if r.HasAuth || r.Public {
			continue
		}
		unguarded++
		siblings, inconsistent := siblingEvidence(r, authed, exposed)
		f := model.Finding{Route: r}
		switch {
		case r.API != "":
			f.Rule, f.Confidence = RuleAPI, model.Medium
			f.Reason = r.API + " is mounted as one handler with no auth on its chain; its operations aren't analysed, so check its own auth (interceptors, security handlers)"
		case r.AuthElsewhere && len(r.PassThrough) == 0 && len(r.CondAuth) == 0:
			// the same registration behind auth on one server or mount and
			// open on another
			f.Rule, f.Confidence = RuleSibling, model.Medium
			f.Reason = "the same route requires auth when reached another way (another server or mount); reached this way it is open"
		case len(r.PassThrough) > 0:
			f.Rule, f.Confidence = RulePassThrough, model.Medium
			f.Reason = strings.Join(r.PassThrough, ", ") + " is named like an auth check, but its code always passes the request on: nothing is rejected"
		case len(r.CondAuth) > 0:
			f.Rule, f.Confidence = RuleConditional, model.Medium
			f.Reason = "auth (" + strings.Join(r.CondAuth, ", ") + ") is only added inside a condition; when it doesn't hold the route is open"
			if strings.Contains(r.CondAuth[0], " [") {
				f.Reason = "auth (" + strings.Join(r.CondAuth, ", ") + "): its own code lets the request through otherwise"
			}
			if isPrivileged(r.Path) {
				f.Reason += " (admin/internal path)"
			}
		case len(r.WeakAuth) > 0:
			// Not escalated by siblings: the weak middleware may well be the
			// app's real check under a name the vocabulary can't confirm.
			f.Rule, f.Confidence = RuleUnverified, model.Medium
			f.Reason = "the only auth-like middleware on the chain (" + strings.Join(r.WeakAuth, ", ") +
				") reads as a loader or a hint, not a check"
			if strings.Contains(r.WeakAuth[0], "[too many branches") {
				f.Reason = "the only auth-like middleware on the chain (" + strings.Join(r.WeakAuth, ", ") +
					") has too many branches to verify that it rejects anything"
			}
			f.Hint = "if it rejects unauthenticated requests, pass -auth " + baseOf(r.WeakAuth[0])
		case inconsistent:
			f.Rule, f.Confidence = RuleSibling, model.High
			f.SiblingsWithAuth = siblings
			f.Reason = "reachable without auth while sibling routes in the same group require it"
			if safeMethod(r.Method) && onlyWrites(siblings) && !isPrivileged(r.Path) {
				// public reads beside authenticated writes is a common design
				f.Confidence = model.Low
				f.Reason = "a public read in a group whose writes require auth (often by design)"
			}
		case pathSiblings(r, authed, exposed) != nil:
			f.Rule, f.Confidence = RuleSibling, model.Medium
			f.SiblingsWithAuth = pathSiblings(r, authed, exposed)
			f.Reason = "the other routes under " + firstSeg(r.Path) + " require auth, on a different router"
		case versionSiblings(r, authed) != nil:
			f.Rule, f.Confidence = RuleSibling, model.Medium
			f.SiblingsWithAuth = versionSiblings(r, authed)
			f.Reason = "the same resource under another API version requires auth"
		case isPrivileged(r.Path):
			f.Rule, f.Confidence = RuleSensitive, model.Medium
			f.Reason = "an admin/internal path with no auth on its chain"
		case appAuth && credentialPath(r.Path):
			f.Rule, f.Confidence = RuleSensitive, model.Medium
			f.Reason = "a token, key or password endpoint with no auth on its chain"
		case appAuth && stateChanging(r.Method) && hasPathParam(r.Path):
			f.Rule, f.Confidence = RuleSensitive, model.Medium
			f.Reason = "a state-changing " + r.Method + " on an object-scoped path with no auth on its chain"
		case appAuth && r.Method == "DELETE":
			f.Rule, f.Confidence = RuleSensitive, model.Medium
			f.Reason = "a DELETE with no auth on its chain"
		default:
			f.Rule, f.Confidence = RuleNoAuth, model.Low
			f.Reason = "reachable without any auth check"
			if why := sensitiveReason(r); why != "" {
				f.Reason += " (" + why + ")"
			}
		}
		if r.HandlerChecks && f.Rule != RuleConditional {
			f.Confidence = lowest(f.Confidence, model.Low)
			f.Reason += "; its handler checks the caller itself (a 401 or identity check in its body)"
		}
		// a token in the URL can stand in for auth on a read; a write keyed by
		// {code} or {fileHash} is still a write
		if p := capabilityParam(r.Path); p != "" && safeMethod(r.Method) && !isPrivileged(r.Path) && f.Confidence.Rank() < model.Low.Rank() {
			f.Confidence = model.Low
			f.Reason += " (the URL carries a " + p + ": a capability link)"
		}
		// a catch-all write (POST /admin/{path...}, DELETE /assets/{path...})
		// is no fallback page
		if isCatchAll(r.Path) && f.Confidence.Rank() < model.Low.Rank() && r.API == "" &&
			!(stateChanging(r.Method) && r.Method != "*") {
			// a bare wildcard is a 404 or SPA fallback far more often than an API
			f.Confidence = model.Low
			f.Reason += " (catch-all route)"
		}
		if r.Partial != "" {
			f.Confidence = downgrade(f.Confidence)
			f.Reason += "; " + r.Partial
		}
		findings = append(findings, f)
	}

	findings = groupSameCause(findings)

	if !appAuth && unguarded > 0 {
		var list []string
		for _, r := range classified {
			if !r.HasAuth && !r.Public {
				list = append(list, fmt.Sprintf("%s %s (%s:%d)", r.Method, r.Path, r.File, r.Line))
			}
		}
		unguarded = len(list)
		findings = append(findings, model.Finding{
			Rule:       RuleNoAuthAtAll,
			Confidence: model.Medium,
			Count:      unguarded,
			Routes:     list,
			Reason: fmt.Sprintf("none of the %d non-public routes has recognised auth middleware; "+
				"either auth happens outside the code (a gateway) or under a name authzmap doesn't know", unguarded),
			Hint: unknownHint(classified, opts),
		})
	}

	sort.SliceStable(findings, func(i, j int) bool {
		a, b := findings[i], findings[j]
		if a.Confidence.Rank() != b.Confidence.Rank() {
			return a.Confidence.Rank() < b.Confidence.Rank()
		}
		if (a.Count > 0) != (b.Count > 0) {
			return a.Count > 0
		}
		if a.Route.Path != b.Route.Path {
			return a.Route.Path < b.Route.Path
		}
		return a.Route.Method < b.Route.Method
	})
	return findings
}

// groupSameCause folds unverified- and conditional-auth findings that share
// the same middleware into one: on a real app that is one decision for the
// user (pass -auth X, or fix the condition), not a hundred.
func groupSameCause(fs []model.Finding) []model.Finding {
	const keepSeparate = 3
	byKey := map[string][]int{}
	for i, f := range fs {
		switch f.Rule {
		case RuleUnverified:
			byKey[f.Rule+"|"+strings.Join(f.Route.WeakAuth, ",")+"|"+string(f.Confidence)] = append(byKey[f.Rule+"|"+strings.Join(f.Route.WeakAuth, ",")+"|"+string(f.Confidence)], i)
		case RuleConditional:
			byKey[f.Rule+"|"+strings.Join(f.Route.CondAuth, ",")+"|"+string(f.Confidence)] = append(byKey[f.Rule+"|"+strings.Join(f.Route.CondAuth, ",")+"|"+string(f.Confidence)], i)
		case RulePassThrough:
			k := f.Rule + "|" + strings.Join(f.Route.PassThrough, ",") + "|" + string(f.Confidence)
			byKey[k] = append(byKey[k], i)
		}
	}
	drop := map[int]bool{}
	for _, idx := range byKey {
		if len(idx) <= keepSeparate {
			continue
		}
		head := &fs[idx[0]]
		head.Count = len(idx)
		for _, i := range idx {
			r := fs[i].Route
			head.Routes = append(head.Routes, fmt.Sprintf("%s %s (%s:%d)", r.Method, r.Path, r.File, r.Line))
			if i != idx[0] {
				drop[i] = true
			}
		}
		switch head.Rule {
		case RulePassThrough:
			head.Reason = fmt.Sprintf("on %d routes %s is named like an auth check, but its code always passes the request on",
				len(idx), strings.Join(head.Route.PassThrough, ", "))
		case RuleUnverified:
			head.Reason = fmt.Sprintf("on %d routes the only auth-like middleware (%s) reads as a loader or a hint, not a check",
				len(idx), strings.Join(head.Route.WeakAuth, ", "))
		default:
			head.Reason = fmt.Sprintf("on %d routes auth (%s) is only added inside a condition; when it doesn't hold they are open",
				len(idx), strings.Join(head.Route.CondAuth, ", "))
		}
	}
	out := fs[:0]
	for i, f := range fs {
		if !drop[i] {
			out = append(out, f)
		}
	}
	return out
}

// localKey groups a route with the siblings that share its scope and its
// first path segment: /me/export beside /me/settings, not beside /{post}.
func localKey(r model.Route) string {
	first := ""
	if segs := staticSegments(r.Path); len(segs) > 0 {
		first = segs[0]
		if isParam(first) {
			first = "{}"
		}
	}
	return r.Scope + " " + first
}

// staticSegments drops leading segments built at runtime (<cfg.BaseURI>),
// which every route of an app shares.
func staticSegments(p string) []string {
	segs := segments(strings.ToLower(p))
	for len(segs) > 0 && strings.HasPrefix(segs[0], "<") {
		segs = segs[1:]
	}
	return segs
}

// pathKey groups a route with every route in the app under the same first
// path segment, whichever router registered it.
func pathKey(r model.Route) string {
	return "path " + firstSeg(r.Path)
}

func firstSeg(p string) string {
	segs := staticSegments(p)
	if len(segs) == 0 || isParam(segs[0]) || isWildcard(segs[0]) {
		return "/"
	}
	return "/" + segs[0]
}

// pathSiblings: the routes under this route's first segment, app-wide, are
// at least 75% guarded and at least three of them. Weaker than a group
// (the routers differ), so it reports at medium.
func pathSiblings(r model.Route, authed map[string][]string, exposed map[string]int) []string {
	k := pathKey(r)
	if firstSeg(r.Path) == "/" {
		return nil
	}
	a := authed[k]
	if len(a) < 3 || len(a)*100 < exposed[k]*minPathPct {
		return nil
	}
	return a
}

// versionKey names a route's resource with the API version wildcarded:
// /api/v2/users/export -> "ver /api/v*/users". "" when there is no version.
func versionKey(p string) string {
	segs := staticSegments(p)
	for i, s := range segs {
		if isVersion(s) {
			end := i + 2
			if end > len(segs) {
				end = len(segs)
			}
			out := append(append([]string(nil), segs[:i]...), "v*")
			out = append(out, segs[i+1:end]...)
			return "ver /" + strings.Join(out, "/")
		}
	}
	return ""
}

func versionOf(p string) string {
	for _, s := range staticSegments(p) {
		if isVersion(s) {
			return s
		}
	}
	return ""
}

// versionSiblings: authed routes for the same resource under a different
// API version (/api/v1/users authed, /api/v2/users not).
func versionSiblings(r model.Route, authed map[string][]string) []string {
	vk := versionKey(r.Path)
	if vk == "" {
		return nil
	}
	mine := versionOf(r.Path)
	var out []string
	for _, id := range authed[vk] {
		_, path, _ := strings.Cut(id, " ")
		if versionOf(path) != mine {
			out = append(out, id)
		}
	}
	return out
}

// siblingEvidence decides whether authed siblings make a gap inconsistent.
// The close siblings (same first segment) must be at least half authed;
// failing that, the whole scope must be three-quarters authed, so a mixed
// group with a handful of checked routes stays quiet. A root router also
// needs two authed routes, since it mixes public and private by design.
func siblingEvidence(r model.Route, authed map[string][]string, exposed map[string]int) ([]string, bool) {
	enough := func(a, e, pct int) bool {
		if a == 0 || a*100 < e*pct {
			return false
		}
		switch r.ScopeKind {
		case model.ScopeGroup, model.ScopePrefix, model.ScopeParam:
			return true
		}
		return a >= 2
	}
	// a parameter as the first segment (/{slug}, /{collection}/...) groups
	// unrelated pages, so only the whole-scope rule applies there
	local := localKey(r)
	if l := authed[local]; !strings.HasSuffix(local, " {}") && enough(len(l), exposed[local], minLocalPct) {
		return l, true
	}
	if s := authed[r.Scope]; enough(len(s), exposed[r.Scope], minScopePct) {
		return s, true
	}
	return nil, false
}

func lowest(a, b model.Confidence) model.Confidence {
	if a.Rank() > b.Rank() {
		return a
	}
	return b
}

func downgrade(c model.Confidence) model.Confidence {
	switch c {
	case model.High:
		return model.Medium
	case model.Medium:
		return model.Low
	}
	return model.Info
}

func baseOf(ref string) string {
	ref = baseName(ref)
	if i := strings.IndexByte(ref, '='); i >= 0 {
		ref = ref[:i]
	}
	return ref
}

// unknownHint lists the most common middleware on unguarded routes that the
// vocabulary doesn't recognise: the likely candidates for -auth.
func unknownHint(routes []model.Route, opts Options) string {
	c := opts.classifier()
	counts := map[string]int{}
	for _, r := range routes {
		if r.HasAuth || r.Public {
			continue
		}
		for _, m := range r.Chain {
			if c.Classify(m) == authsig.Unknown {
				counts[baseOf(m)]++
			}
		}
	}
	type kv struct {
		k string
		v int
	}
	var top []kv
	for k, v := range counts {
		top = append(top, kv{k, v})
	}
	sort.Slice(top, func(i, j int) bool {
		if top[i].v != top[j].v {
			return top[i].v > top[j].v
		}
		return top[i].k < top[j].k
	})
	if len(top) == 0 {
		return "if a gateway or proxy authenticates every request, pass -allow \"/**\""
	}
	var names []string
	for i := 0; i < len(top) && i < 4; i++ {
		names = append(names, top[i].k)
	}
	return "if one of these is your auth middleware, pass -auth NAME: " + strings.Join(names, ", ")
}

func safeMethod(m string) bool { return m == "GET" || m == "HEAD" }

// capabilityParam returns the name of a path parameter that is itself the
// credential (/t/:token, /invite/{code}), or "".
func capabilityParam(p string) string {
	for _, s := range segments(strings.ToLower(p)) {
		if !isParam(s) {
			continue
		}
		name := strings.Trim(s, "{}:<>*")
		if i := strings.IndexByte(name, ':'); i >= 0 {
			name = name[:i] // gorilla {code:[a-z]+}
		}
		for _, w := range []string{"token", "secret", "signature", "sig", "hash", "invite", "code"} {
			if strings.Contains(name, w) {
				return name
			}
		}
	}
	return ""
}

// isCatchAll: "/*", "/api/*", "/{path...}" with no parameter before the
// wildcard.
func isCatchAll(p string) bool {
	segs := segments(p)
	if len(segs) == 0 || len(segs) > 2 || !isWildcard(segs[len(segs)-1]) {
		return false
	}
	return len(segs) == 1 || !isParam(segs[0])
}

// onlyWrites: every "METHOD /path" in siblings uses a state-changing method.
func onlyWrites(siblings []string) bool {
	for _, s := range siblings {
		m, _, _ := strings.Cut(s, " ")
		if !stateChanging(m) || m == "*" {
			return false
		}
	}
	return true
}

func stateChanging(m string) bool {
	switch m {
	case "POST", "PUT", "DELETE", "PATCH":
		return true
	}
	return false
}

func isPrivileged(p string) bool {
	// /api/%61dmin is /api/admin once the server decodes it
	if d, err := url.PathUnescape(p); err == nil {
		p = d
	}
	for _, s := range segments(strings.ToLower(p)) {
		if privileged[s] || s == "pprof" {
			return true
		}
	}
	return false
}

var credentialSeg = set("token", "tokens", "apikey", "apikeys", "api-key", "api-keys", "api_key", "api_keys",
	"keys", "secret", "secrets", "password", "passwords", "credential", "credentials", "revoke", "revoke-all")

// credentialPath: a route that hands out, changes or revokes credentials.
func credentialPath(p string) bool {
	for _, s := range segments(strings.ToLower(p)) {
		if credentialSeg[s] {
			return true
		}
	}
	return false
}

func sensitiveReason(r model.Route) string {
	var parts []string
	if stateChanging(r.Method) && r.Method != "*" {
		parts = append(parts, "state-changing "+r.Method)
	}
	if hasPathParam(r.Path) {
		parts = append(parts, "object-scoped path")
	}
	return strings.Join(parts, ", ")
}

func hasPathParam(path string) bool {
	for _, seg := range strings.Split(path, "/") {
		if strings.HasPrefix(seg, "{") || strings.HasPrefix(seg, ":") {
			return true
		}
	}
	return false
}

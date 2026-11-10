package analyze

import (
	"strings"
	"testing"

	"github.com/NotChaosuu/authzmap/pkg/model"
	"github.com/NotChaosuu/authzmap/pkg/scan"
)

func route(method, path, scope string, chain ...string) model.Route {
	if chain == nil {
		chain = []string{}
	}
	return model.Route{Method: method, Path: path, Scope: scope, ScopeKind: model.ScopeGroup, Chain: chain}
}

func root(r model.Route) model.Route { r.ScopeKind = model.ScopeRoot; return r }

func find(fs []model.Finding, method, path string) (model.Finding, bool) {
	for _, f := range fs {
		if f.Route.Method == method && f.Route.Path == path {
			return f, true
		}
	}
	return model.Finding{}, false
}

func expectConf(t *testing.T, fs []model.Finding, method, path string, want model.Confidence) model.Finding {
	t.Helper()
	f, ok := find(fs, method, path)
	if !ok {
		if want != "" {
			t.Errorf("%s %s: no finding, want %s", method, path, want)
		}
		return f
	}
	if want == "" {
		t.Errorf("%s %s: unexpected %s %s finding", method, path, f.Confidence, f.Rule)
	} else if f.Confidence != want {
		t.Errorf("%s %s: %s %s (%s), want %s", method, path, f.Confidence, f.Rule, f.Reason, want)
	}
	return f
}

func TestSiblingInconsistencyIsHigh(t *testing.T) {
	fs := Findings([]model.Route{
		route("GET", "/admin/users", "g1", "RequireAuth"),
		route("GET", "/admin/stats", "g1"),
	}, Options{})
	f := expectConf(t, fs, "GET", "/admin/stats", model.High)
	if f.Rule != RuleSibling || len(f.SiblingsWithAuth) != 1 {
		t.Errorf("%+v", f)
	}
	expectConf(t, fs, "GET", "/admin/users", "")
}

// On a root router one authed route doesn't make every public route
// a high-confidence gap (gotify's login, CORS preflight, version, info).
func TestRootRouterNeedsMajority(t *testing.T) {
	fs := Findings([]model.Route{
		root(route("GET", "/plugin", "root", "RequireClient")),
		root(route("GET", "/gotifyinfo", "root")),
		root(route("GET", "/about", "root")),
		root(route("GET", "/stream", "root")),
	}, Options{})
	for _, f := range fs {
		if f.Confidence == model.High {
			t.Errorf("root-router minority produced HIGH: %s %s", f.Route.Method, f.Route.Path)
		}
	}
}

// A mixed group with a few checked routes stays quiet; the close siblings
// (same first segment) carry the signal.
func TestLocalSiblings(t *testing.T) {
	fs := Findings([]model.Route{
		route("GET", "/me/settings", "g", "RequireUser"),
		route("GET", "/me/export", "g"),
		route("GET", "/{post}", "g"),
		route("GET", "/about", "g"),
		route("GET", "/feed", "g"),
	}, Options{})
	expectConf(t, fs, "GET", "/me/export", model.High)
	expectConf(t, fs, "GET", "/{post}", model.Low)
}

// The other routes under /application are authed on another router.
func TestPathSiblings(t *testing.T) {
	fs := Findings([]model.Route{
		route("GET", "/application", "a", "RequireClient"),
		route("POST", "/application", "a", "RequireClient"),
		route("PUT", "/application/:id", "a", "RequireClient"),
		root(route("GET", "/application/export", "root")),
	}, Options{})
	f := expectConf(t, fs, "GET", "/application/export", model.Medium)
	if f.Rule != RuleSibling {
		t.Errorf("rule = %s", f.Rule)
	}
}

func TestPublicReadBesideAuthedWrites(t *testing.T) {
	fs := Findings([]model.Route{
		route("POST", "/api/gists", "g", "apiRequireAuth"),
		route("DELETE", "/api/gists/:id", "g", "apiRequireAuth"),
		route("GET", "/api/gists", "g"),
	}, Options{})
	expectConf(t, fs, "GET", "/api/gists", model.Low)
}

// State-changing alone isn't MEDIUM, and an app with no auth at all
// gets one summary instead of a finding per route.
func TestSensitiveNeedsTwoSignals(t *testing.T) {
	authed := []model.Route{
		route("GET", "/me", "g0", "RequireAuth"),
		route("DELETE", "/users/{id}", "g1"),
		route("POST", "/feedback", "g2"),
		route("GET", "/s/{id}", "g3"),
		route("GET", "/internal/debug", "g4"),
	}
	fs := Findings(authed, Options{})
	expectConf(t, fs, "DELETE", "/users/{id}", model.Medium)
	expectConf(t, fs, "POST", "/feedback", model.Low)
	expectConf(t, fs, "GET", "/s/{id}", model.Low)
	expectConf(t, fs, "GET", "/internal/debug", model.Medium)

	none := Findings([]model.Route{route("DELETE", "/users/{id}", "g"), route("GET", "/about", "g")}, Options{})
	expectConf(t, none, "DELETE", "/users/{id}", model.Low)
	summary := false
	for _, f := range none {
		if f.Rule == RuleNoAuthAtAll && f.Count == 2 && f.Confidence == model.Medium {
			summary = true
		}
	}
	if !summary {
		t.Errorf("no-auth-anywhere summary missing: %+v", none)
	}
}

// Lookalike middleware (tokenBucket, quota.Enforce) doesn't hide /admin routes.
func TestLookalikeMiddlewareStillReported(t *testing.T) {
	fs := Findings([]model.Route{
		root(route("GET", "/admin/real", "r", "RequireAuth")),
		root(route("GET", "/admin/ratelimited", "r", "tokenBucket")),
		root(route("GET", "/admin/sessionload", "r", "session.Load")),
		root(route("GET", "/admin/quota", "r", "quota.Enforce")),
		root(route("GET", "/admin/oauthclient", "r", "oauthClient.Attach")),
		root(route("GET", "/admin/feature", "r", "featureGuard")),
	}, Options{})
	for _, p := range []string{"/admin/ratelimited", "/admin/sessionload", "/admin/quota", "/admin/oauthclient", "/admin/feature"} {
		expectConf(t, fs, "GET", p, model.Medium)
	}
	if f, _ := find(fs, "GET", "/admin/sessionload"); f.Rule != RuleUnverified || f.Hint == "" {
		t.Errorf("weak middleware should be unverified-auth with a hint: %+v", f)
	}
}

func TestUnverifiedIsGrouped(t *testing.T) {
	var rs []model.Route
	for _, p := range []string{"/a", "/b", "/c", "/d", "/e"} {
		rs = append(rs, route("GET", p, "g", "session.Load"))
	}
	grouped := 0
	for _, f := range Findings(rs, Options{}) {
		switch f.Rule {
		case RuleUnverified:
			if f.Count != 5 || len(f.Routes) != 5 {
				t.Errorf("grouped finding = %+v", f)
			}
			grouped++
		case RuleNoAuthAtAll:
		default:
			t.Errorf("unexpected %+v", f)
		}
	}
	if grouped != 1 {
		t.Errorf("want exactly one grouped finding, got %d", grouped)
	}
}

func TestConditionalAuth(t *testing.T) {
	r := route("GET", "/debug/pprof/profile", "g", "auth=basicAuth")
	r.Conditional = []string{"auth=basicAuth"}
	fs := Findings([]model.Route{r, route("GET", "/x", "h", "RequireAuth")}, Options{})
	f := expectConf(t, fs, "GET", "/debug/pprof/profile", model.Medium)
	if f.Rule != RuleConditional {
		t.Errorf("rule = %s", f.Rule)
	}
}

func TestBodyVerifiedMiddleware(t *testing.T) {
	r := route("GET", "/feeds", "g", "webSession.handle")
	r.Enforcing = []string{"webSession.handle"}
	if cl := Classify([]model.Route{r}, Options{}); !cl[0].HasAuth {
		t.Error("middleware whose body rejects anonymous requests counts as auth")
	}
	r.Chain, r.Enforcing = []string{"OptionalAuth"}, []string{"OptionalAuth"}
	if cl := Classify([]model.Route{r}, Options{}); cl[0].HasAuth {
		t.Error("an optional-auth name outranks a 401 in its body")
	}
}

func TestHandlerChecksLowersConfidence(t *testing.T) {
	r := route("DELETE", "/albums/:uid", "g")
	r.HandlerChecks = true
	fs := Findings([]model.Route{r, route("GET", "/me", "h", "RequireAuth")}, Options{})
	expectConf(t, fs, "DELETE", "/albums/:uid", model.Low)
}

func TestPartialRouteDowngraded(t *testing.T) {
	r := route("DELETE", "/users/{id}", "g")
	r.Partial = "router arrives as parameter r"
	fs := Findings([]model.Route{r, route("GET", "/me", "h", "RequireAuth")}, Options{})
	expectConf(t, fs, "DELETE", "/users/{id}", model.Low)
}

// Entrypoints under a prefix, OPTIONS, probes under a
// neutral prefix, docs and capability links.
func TestPublicPaths(t *testing.T) {
	public := map[string]string{
		"POST /auth/local/login":                         "",
		"POST /auth/oidc/external/token":                 "",
		"POST /auth/oidc/external/authorize":             "",
		"GET /auth/oidc/login":                           "",
		"GET /api/v1/auth/openid/:provider/callback":     "",
		"OPTIONS /*any":                                  "",
		"GET /version":                                   "",
		"GET /system/ping":                               "",
		"GET /api/health":                                "",
		"GET /v1/metrics":                                "",
		"GET /api/openapi.yaml":                          "",
		"GET /oauth/:provider":                           "",
		"POST /webauthn/login/finish":                    "",
		"GET <conf.BaseUri(\"/.well-known/jwks.json\")>": "",
		"GET /static/js/app.js":                          "",
		"POST /webhooks/stripe":                          "",
	}
	private := []string{
		"GET /auth/oidc/elevate", "POST /settings/token", "GET /users/{id}/status",
		"GET /oauth/clients/new", "GET /admin/users", "DELETE /docs/{id}", "GET /users/{id}/token",
		"GET /admin/login", "DELETE /assets/{id}", "POST /api/nodes/register", "POST /password/change",
		"DELETE /webhooks/{id}", "POST /api/public/x", "POST /admin/jobs/callback",
	}
	for k := range public {
		m, p, _ := strings.Cut(k, " ")
		if !isPublic(m, p, Options{}) {
			t.Errorf("%s should be public", k)
		}
	}
	for _, k := range private {
		m, p, _ := strings.Cut(k, " ")
		if isPublic(m, p, Options{}) {
			t.Errorf("%s should not be public", k)
		}
	}
}

func TestCapabilityAndCatchAllAreLow(t *testing.T) {
	fs := Findings([]model.Route{
		route("GET", "/albums/:uid/t/:token/:size", "g"),
		route("GET", "/albums/:uid", "g", "RequireAuth"),
		root(route("*", "/*", "r")),
		root(route("GET", "/x", "r", "RequireAuth")),
		root(route("GET", "/y", "r", "RequireAuth")),
	}, Options{})
	expectConf(t, fs, "GET", "/albums/:uid/t/:token/:size", model.Low)
	expectConf(t, fs, "*", "/*", model.Low)
}

func TestAllowGlobs(t *testing.T) {
	routes := []model.Route{route("GET", "/internal/debug/x", "g")}
	if fs := Findings(routes, Options{AllowGlobs: []string{"/internal/**"}}); len(fs) != 0 {
		t.Errorf("** should cover any depth: %+v", fs)
	}
	if !matchGlob("/a/*", "/a/b") || matchGlob("/a/*", "/a/b/c") || !matchGlob("/**", "/anything/at/all") {
		t.Error("glob semantics")
	}
}

func TestNoDefaultAllow(t *testing.T) {
	fs := Findings([]model.Route{route("POST", "/login", "g")}, Options{NoDefaultAllow: true})
	if len(fs) == 0 {
		t.Error("-no-default-allow must report /login")
	}
}

// End to end on the chi fixture: the only gap is the admin route.
func TestFixtureEndToEnd(t *testing.T) {
	res, err := scan.Scan("../scan/testdata/chi_basic")
	if err != nil {
		t.Fatal(err)
	}
	var high []model.Finding
	for _, f := range Findings(res.Routes, Options{}) {
		if f.Confidence == model.High || f.Confidence == model.Medium {
			high = append(high, f)
		}
	}
	if len(high) != 1 || high[0].Route.Path != "/admin/stats" || high[0].Confidence != model.High {
		t.Fatalf("want high on /admin/stats only, got %+v", high)
	}
}

func TestPublicTrees(t *testing.T) {
	for _, p := range []string{"/api/public/lists", "/admin/static/*", "/api/v1/info", "/api/v1/docs/redoc.standalone.js"} {
		if !isPublic("GET", p, Options{}) {
			t.Errorf("%s should be public", p)
		}
	}
	for _, p := range []string{"/users/{id}/info", "/admin/static", "/admin/users"} {
		if isPublic("GET", p, Options{}) {
			t.Errorf("%s should not be public", p)
		}
	}
}

// A registration reached through two chains is as protected as the weaker
// one: AuthA on one server and AuthB on the other is authed; auth on only
// one of them isn't.
func TestVariantsJudgedByWeakest(t *testing.T) {
	two := route("DELETE", "/files/{id}", "s")
	two.Variants = []model.Variant{{Chain: []string{"RequireAuth"}}, {Chain: []string{"TokenAuth"}}}
	one := route("DELETE", "/blobs/{id}", "s")
	one.Variants = []model.Variant{{Chain: []string{"RequireAuth"}}, {Chain: []string{"logger"}}}
	got := Classify([]model.Route{two, one}, Options{})
	if !got[0].HasAuth {
		t.Error("two auths: want authed")
	}
	if got[1].HasAuth {
		t.Error("auth on one path only: want open")
	}
}

// A path the server decodes (%61dmin) and credential endpoints are
// sensitive; a catch-all write is not a fallback page.
func TestSensitivePaths(t *testing.T) {
	authed := route("GET", "/api/me", "s", "RequireAuth")
	fs := Findings([]model.Route{authed,
		route("GET", "/api/%61dmin/users", "o"),
		route("POST", "/oauth/token/revoke-all", "o"),
		route("POST", "/admin/{path...}", "o"),
	}, Options{})
	expectConf(t, fs, "GET", "/api/%61dmin/users", model.Medium)
	expectConf(t, fs, "POST", "/oauth/token/revoke-all", model.Medium)
	expectConf(t, fs, "POST", "/admin/{path...}", model.Medium)
}

// A route open on one server or mount and authed on another is reported
// as such, and the app isn't summarised as having no auth at all.
func TestAuthElsewhere(t *testing.T) {
	r := route("DELETE", "/users/{id}", "s")
	r.Variants = []model.Variant{{Chain: []string{"RequireAuth"}}, {Chain: []string{}}}
	fs := Findings([]model.Route{r}, Options{})
	f := expectConf(t, fs, "DELETE", "/users/{id}", model.Medium)
	if f.Rule != RuleSibling {
		t.Errorf("rule = %s", f.Rule)
	}
	for _, f := range fs {
		if f.Rule == RuleNoAuthAtAll {
			t.Error("no-auth-anywhere while one server has auth")
		}
	}
}

// Middleware whose code has too many branches to read is unverified, not
// trusted by its name.
func TestTooBranchyIsUnverified(t *testing.T) {
	m := "RequireAuth [too many branches to verify]"
	r := route("DELETE", "/users/{id}", "s", m)
	r.Unverified = []string{m}
	f := expectConf(t, Findings([]model.Route{r}, Options{}), "DELETE", "/users/{id}", model.Medium)
	if f.Rule != RuleUnverified || !strings.Contains(f.Reason, "too many branches") {
		t.Errorf("rule=%s reason=%q", f.Rule, f.Reason)
	}
}

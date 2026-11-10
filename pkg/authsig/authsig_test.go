package authsig

import (
	"strings"
	"testing"
)

func expect(t *testing.T, want Kind, names ...string) {
	t.Helper()
	for _, n := range names {
		if got := ClassifyName(n); got != want {
			t.Errorf("ClassifyName(%q) = %v, want %v", n, got, want)
		}
	}
}

func TestAuthNames(t *testing.T) {
	expect(t, Auth,
		"AuthMiddleware", "middleware.RequireAuth", "JWTAuth", "RequireRole",
		"AuthRequired()", "jwtauth.Authenticator", "SessionAuth", "EnsureLoggedIn",
		"middleware.Authenticate", "RBACGuard", "casbin.Enforce", "AdminAuth()",
		"requireLogin", "BearerToken", "authentication.RequireClient", "echojwt.WithConfig",
		"keyauth.New", "middleware.BasicAuth", "metricsBasicAuth", "logged", "isAdmin",
		"requireSession", "withApiKeyAuth", "adminGuard", "authMW=jwtauth.Authenticator")
}

// Real non-auth middleware, much of it named with an auth word (tokenBucket,
// quota.Enforce, authLogger). A substring match would count those as auth
// and hide the route from every finding.
func TestNonAuthLookalikes(t *testing.T) {
	expect(t, NotAuth, "tokenBucket", "quota.Enforce", "featureGuard", "csrf.Protect",
		"middleware.Logger", "CORS", "Gzip", "RateLimiter", "Timeout", "Metrics",
		"AuthorMiddleware", "Authors", "authLogger", "loginRateLimiter", "accessLog")
	expect(t, Weak, "session.Load", "oauthClient.Attach", "OAuth2Middleware", "token.Refresh",
		"jwtauth.Verifier", "authMW=jwtauth.Verifier") // Verifier parses; Authenticator rejects
	for _, n := range []string{"tokenBucket", "session.Load", "quota.Enforce", "oauthClient.Attach", "featureGuard"} {
		if ClassifyName(n) == Auth {
			t.Errorf("%q must not count as auth", n)
		}
	}
}

// Middleware that lets anonymous requests through is not enforcement.
func TestOptionalAuthIsWeak(t *testing.T) {
	expect(t, Weak, "OptionalAuth", "auth.Optional", "NoAuth", "SkipAuth", "MaybeAuthenticated",
		"authentication.OptionalAdmin", "handler.Web(UserLevelNone)")
}

func TestUnknownNames(t *testing.T) {
	expect(t, Unknown, "FooBar", "CustomThing", "myMiddleware", "RequestID", "currentUser")
}

func TestOverrides(t *testing.T) {
	c := NewClassifier([]string{"webSessionMiddleware.handle", "gate"}, []string{"JWTAuth"})
	if c.Classify("webSessionMiddleware.handle=newWebSessionMiddleware.handle") != Auth {
		t.Error("-auth must match the variable name of an assigned middleware")
	}
	if c.Classify("mw.gate") != Auth {
		t.Error("-auth must match the last element")
	}
	if c.Classify("JWTAuth") != NotAuth {
		t.Error("-not-auth must win over the vocabulary")
	}
	if !c.Optional("OptionalAuth") || c.Optional("RequireAuth") {
		t.Error("Optional misreads names")
	}
}

func TestWords(t *testing.T) {
	cases := map[string]string{
		"JWTAuth":                "jwt auth",
		"OAuth2Middleware":       "oauth2 middleware",
		"middleware.RequireAuth": "middleware require auth",
		"require_login":          "require login",
		"APIKeyAuth":             "apikey auth",
	}
	for in, want := range cases {
		if got := strings.Join(words(in), " "); got != want {
			t.Errorf("words(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAnyAuth(t *testing.T) {
	if !AnyAuth([]string{"Logger", "RequireAuth", "CORS"}) || AnyAuth([]string{"Logger", "CORS", "session.Load"}) {
		t.Error("AnyAuth")
	}
}

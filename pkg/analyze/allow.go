package analyze

import (
	"regexp"
	"strings"
)

// Public by design: auth entrypoints, health probes, webhook receivers (they
// authenticate by signature), static assets and API docs. Matched by path
// segment so /api/v1/auth/login counts as well as /login, but never under an
// admin or internal area, never for an object-scoped path (/users/{id}/token),
// and only for the methods such a route serves: a DELETE on an asset tree or
// a POST to a docs page is not public.

var (
	// last segment of an auth entrypoint
	sessionEntry = set("login", "logout", "signin", "signout", "sign-in", "sign-out")
	accountEntry = set("signup", "sign-up", "register", "registration")
	flowEntry    = set("callback", "oauth2callback", "authorize")
	authEntry    = set("login", "logout", "signin", "signout", "sign-in", "sign-out", "signup", "sign-up",
		"register", "registration", "callback", "oauth2callback", "authorize", "verify-email", "confirm-email")
	// where a registration endpoint registers a person
	accountArea = set("user", "users", "account", "accounts", "member", "members", "auth", "oauth", "oidc", "sso")
	resetPaths  = set("forgot", "forgot-password", "forgotpassword", "password/forgot", "password/reset",
		"reset-password", "resetpassword", "password-reset", "password/forgot-password", "password/reset-password")
	hookActions = set("delete", "remove", "update", "edit", "reset", "rotate", "disable", "enable", "config", "settings", "secret")
	// token endpoints: only under an auth area or at the top level, since
	// /settings/token usually mints a token for a signed-in user
	tokenEntry = set("token", "refresh", "refresh-token", "revoke", "introspect", "jwks", "jwks.json")
	authArea   = set("auth", "oauth", "oauth2", "oidc", "openid", "openid-connect", "sso", "saml", "session",
		"sessions", "login", "connect", "external", "local", "provider", "providers")
	// health and info probes, at the top level or under a neutral prefix
	probe = set("health", "healthz", "healthcheck", "livez", "readyz", "ready", "live", "liveness",
		"readiness", "ping", "status", "metrics", "version", "info", "docs", "docs.json", "redoc.standalone.js")
	// prefixes that don't change what a probe means
	neutral = set("api", "system", "sys", "ops", "_", "-", "docs")
	// first segment of public asset and docs trees
	assetRoot = set("static", "assets", "public", "dist", "swagger", "swagger-ui", "openapi", "redoc", "docs", "api-docs", ".well-known")
	// webhook receivers
	hookRoot  = set("webhook", "webhooks", "hooks", "callbacks")
	rootFiles = set("favicon.ico", "robots.txt", "sitemap.xml", "manifest.json", "site.webmanifest", "apple-touch-icon.png")
	// steps of a login flow under an entry segment: /webauthn/login/finish
	flowStep = set("begin", "start", "finish", "complete", "verify", "init", "challenge", "options",
		"confirm", "otp", "mfa", "totp", "assertion", "redirect", "consent")
	// API description files, wherever they're served
	docFiles = set("openapi.yaml", "openapi.yml", "openapi.json", "swagger.json", "swagger.yaml", "swagger.yml")
	// first segment of an OAuth/SSO entry tree: /oauth/{provider}
	oauthRoot  = set("oauth", "oauth2", "oidc", "openid", "sso", "saml")
	privileged = set("admin", "admins", "internal", "manage", "management", "debug", "superuser", "staff")
)

// isPublic also reads a path built at runtime by the literal inside it:
// <conf.BaseUri("/.well-known/jwks.json")> is judged as /.well-known/jwks.json.
func isPublic(method, path string, opts Options) bool {
	if publicPath(method, path, opts) {
		return true
	}
	if lit := literalPath(path); lit != "" && lit != path {
		return publicPath(method, lit, opts)
	}
	return false
}

var (
	dynSpan = regexp.MustCompile(`<[^>]*>`)
	quoted  = regexp.MustCompile(`"[^"]*"`)
)

func literalPath(p string) string {
	if !strings.Contains(p, "<") {
		return ""
	}
	out := dynSpan.ReplaceAllStringFunc(p, func(span string) string {
		var b strings.Builder
		for _, q := range quoted.FindAllString(span, -1) {
			b.WriteString(strings.Trim(q, `"`))
		}
		return b.String()
	})
	return "/" + strings.TrimLeft(strings.ReplaceAll(out, "//", "/"), "/")
}

func publicPath(method, path string, opts Options) bool {
	p := strings.ToLower(strings.TrimRight(path, "/"))
	if p == "" {
		p = "/"
	}
	for _, g := range opts.AllowGlobs {
		if matchGlob(g, p) {
			return true
		}
	}
	if opts.NoDefaultAllow {
		return false
	}
	if method == "OPTIONS" {
		return true // CORS preflight carries no credentials by design
	}
	segs := segments(p)
	if len(segs) == 0 {
		return false
	}
	read := method == "GET" || method == "HEAD" || method == "*"
	last := segs[len(segs)-1]
	// an admin UI's own static files (GET /admin/static/*) are still files
	for i, s := range segs {
		if read && (s == "static" || s == "assets") && i+1 < len(segs) && isWildcard(last) {
			return true
		}
	}
	for _, s := range segs {
		if privileged[s] {
			return false // nothing else under an admin or internal area is public by name
		}
	}
	if len(segs) == 1 && rootFiles[segs[0]] {
		return read
	}
	// asset and docs trees serve files: reads only, and docs/{id} is content
	if assetRoot[segs[0]] {
		return read && (segs[0] != "docs" || !anyParam(segs))
	}
	for i, s := range segs {
		if (s == "static" || s == "assets" || s == "dist") && i+1 < len(segs) && isWildcard(last) {
			return read
		}
		// /api/public/lists: an app that names a read tree public means it
		if s == "public" && i > 0 && i < len(segs)-1 {
			return read && !anyParam(segs[:i])
		}
	}
	if docFiles[last] {
		return read && !anyParam(segs)
	}
	// webhook receivers authenticate by signature; managing them doesn't
	if hookRoot[segs[0]] || strings.HasPrefix(segs[0], "webhook") {
		return webhookReceiver(method, segs)
	}
	if resetFlow(segs) {
		return method == "GET" || method == "POST" || method == "*"
	}
	entryMethod := method == "GET" || method == "POST" || method == "*"
	// /oauth/{provider}, /oauth/{provider}/callback
	if oauthRoot[segs[0]] && len(segs) <= 3 && entryMethod {
		ok := true
		for _, s := range segs[1:] {
			if !isParam(s) && !authEntry[s] && !flowStep[s] {
				ok = false
			}
		}
		if ok {
			return true
		}
	}
	// /webauthn/login/finish: a flow step right under an entry segment
	if len(segs) >= 2 && flowStep[last] && authEntry[segs[len(segs)-2]] {
		return entryMethod
	}
	if isWildcard(last) && len(segs) > 1 {
		last = segs[len(segs)-2]
		segs = segs[:len(segs)-1]
	}
	head := segs[:len(segs)-1]
	if anyParam(segs) {
		// /auth/openid/{provider}/callback: a provider id is fine under an auth area
		return entryMethod && authEntry[last] && allIn(head, authArea, true) && hasAny(head, authArea)
	}
	switch {
	case sessionEntry[last]:
		return entryMethod || (method == "DELETE" && (last == "logout" || last == "signout" || last == "sign-out"))
	case accountEntry[last]:
		// registering a user or an account, not a node or a device
		return entryMethod && allIn(head, accountArea, false)
	case flowEntry[last]:
		// a callback or authorize step belongs to an auth flow
		return entryMethod && (allIn(head, neutral, false) || hasAny(head, authArea))
	case authEntry[last]:
		return entryMethod
	case tokenEntry[last]:
		return (method == "POST" || method == "*" || read && strings.HasPrefix(last, "jwks")) && allIn(head, authArea, true)
	case probe[last]:
		return read && allIn(head, neutral, false)
	}
	return false
}

// webhookReceiver: POST (or a GET verification) to a named receiver like
// /webhooks/stripe or /hooks/{provider}; never a DELETE, an id-scoped action
// (/hooks/{id}/delete) or a nested management path.
func webhookReceiver(method string, segs []string) bool {
	if method != "POST" && method != "GET" && method != "*" {
		return false
	}
	if len(segs) > 3 {
		return false
	}
	for i, s := range segs[1:] {
		if isParam(s) && i+1 < len(segs)-1 {
			return false // a parameter followed by more path is an action on a hook
		}
		if hookActions[s] {
			return false
		}
	}
	return true
}

// resetFlow: the forgot/reset-password pages themselves, optionally under a
// neutral or auth prefix and with a trailing token, not anything that merely
// starts with "/forgot" or "/password/".
func resetFlow(segs []string) bool {
	rest := segs
	for len(rest) > 0 && (neutral[rest[0]] || isVersion(rest[0]) || authArea[rest[0]] || rest[0] == "user" || rest[0] == "users" || rest[0] == "account") {
		rest = rest[1:]
	}
	if len(rest) > 0 && isParam(rest[len(rest)-1]) {
		rest = rest[:len(rest)-1]
	}
	return resetPaths[strings.Join(rest, "/")]
}

func hasAny(segs []string, in map[string]bool) bool {
	for _, s := range segs {
		if in[s] {
			return true
		}
	}
	return false
}

// segments splits a path on "/", keeping a runtime-built span such as
// <cfg.Base + "/api/v1"> whole: the slashes inside it aren't path structure
// we can see.
func segments(p string) []string {
	var out []string
	var cur strings.Builder
	depth := 0
	flush := func() {
		if cur.Len() > 0 {
			out = append(out, cur.String())
			cur.Reset()
		}
	}
	for _, r := range p {
		switch {
		case r == '<':
			depth++
		case r == '>' && depth > 0:
			depth--
		case r == '/' && depth == 0:
			flush()
			continue
		}
		cur.WriteRune(r)
	}
	flush()
	return out
}

func isParam(s string) bool {
	return strings.HasPrefix(s, "{") || strings.HasPrefix(s, ":") || strings.HasPrefix(s, "<")
}

func isWildcard(s string) bool {
	return s == "*" || strings.HasPrefix(s, "*") || strings.HasSuffix(s, "...}")
}

func anyParam(segs []string) bool {
	for _, s := range segs {
		if isParam(s) || isWildcard(s) {
			return true
		}
	}
	return false
}

// allIn reports whether every segment is in the set or neutral (api, v1).
// With oneIsEnough, any segments may appear as long as one is in the set
// (/auth/oidc/external/token).
func allIn(segs []string, in map[string]bool, oneIsEnough bool) bool {
	hit, other := false, false
	for _, s := range segs {
		switch {
		case in[s]:
			hit = true
		case isVersion(s), neutral[s]:
		default:
			other = true
		}
	}
	if !other {
		return true
	}
	return oneIsEnough && hit
}

func isVersion(s string) bool {
	if len(s) < 2 || s[0] != 'v' {
		return false
	}
	for _, r := range s[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// matchGlob supports a trailing * (one segment) or ** (any depth); anything
// else is an exact match.
func matchGlob(pattern, path string) bool {
	pattern = strings.ToLower(strings.TrimRight(strings.TrimSpace(pattern), "/"))
	if pattern == "" {
		return false
	}
	switch {
	case strings.HasSuffix(pattern, "/**"):
		base := strings.TrimSuffix(pattern, "/**")
		return path == base || strings.HasPrefix(path, base+"/") || base == ""
	case strings.HasSuffix(pattern, "**"):
		return strings.HasPrefix(path, strings.TrimSuffix(pattern, "**"))
	case strings.HasSuffix(pattern, "/*"):
		base := strings.TrimSuffix(pattern, "/*")
		if !strings.HasPrefix(path, base+"/") {
			return path == base
		}
		return !strings.Contains(path[len(base)+1:], "/")
	case strings.HasSuffix(pattern, "*"):
		return strings.HasPrefix(path, strings.TrimSuffix(pattern, "*"))
	}
	return path == pattern
}

func set(xs ...string) map[string]bool {
	m := make(map[string]bool, len(xs))
	for _, x := range xs {
		m[x] = true
	}
	return m
}

// Package authsig decides from a middleware's name whether it enforces
// authentication or authorization.
//
// Names are split into words (camelCase, snake_case, dotted) and judged on
// whole words, so tokenBucket isn't "token" auth and Authors isn't "auth".
// A name that only hints at auth (session.Load, featureGuard, OAuth2Client)
// is Weak: reported as unverified rather than trusted.
package authsig

import (
	"strings"
	"unicode"
)

// Kind is the verdict on a middleware name.
type Kind int

const (
	Unknown Kind = iota // the name says nothing either way
	Auth                // the name says it enforces auth
	Weak                // the name hints at auth; not trusted
	NotAuth             // the name says it's something else: Logger, RateLimit
)

func (k Kind) String() string {
	switch k {
	case Auth:
		return "auth"
	case Weak:
		return "weak"
	case NotAuth:
		return "not-auth"
	default:
		return "unknown"
	}
}

// Classifier holds user overrides on top of the built-in vocabulary.
type Classifier struct {
	auth    []string
	notAuth []string
}

// NewClassifier takes the names passed with -auth and -not-auth. A name
// matches the full reference (mw.RequireUser) or its last element.
func NewClassifier(auth, notAuth []string) *Classifier {
	return &Classifier{auth: lowerAll(auth), notAuth: lowerAll(notAuth)}
}

func lowerAll(xs []string) []string {
	var out []string
	for _, x := range xs {
		if x = strings.ToLower(strings.TrimSpace(x)); x != "" {
			out = append(out, x)
		}
	}
	return out
}

var defaultClassifier = &Classifier{}

// ClassifyName classifies with the built-in vocabulary only.
func ClassifyName(ref string) Kind { return defaultClassifier.Classify(ref) }

// AnyAuth reports whether any name in the chain enforces auth.
func AnyAuth(names []string) bool { return defaultClassifier.AnyAuth(names) }

// AnyAuth reports whether any name in the chain enforces auth, with the
// user's overrides applied.
func (c *Classifier) AnyAuth(names []string) bool {
	for _, n := range names {
		if c.Classify(n) == Auth {
			return true
		}
	}
	return false
}

// Classify judges one chain entry. The scanner renders a middleware held in
// a variable as "authMW=jwtauth.Verifier"; the part after '=' is what it was
// built from and is judged first.
func (c *Classifier) Classify(ref string) Kind {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return Unknown
	}
	if k, ok := c.override(ref); ok {
		return k
	}
	// chi's jwtauth.Verifier only parses the token into the context; the
	// rejecting half is jwtauth.Authenticator
	if strings.Contains(strings.ToLower(ref), "jwtauth.verif") {
		return Weak
	}
	// authboss.Middleware2(ab, authboss.RequireNone, ...) always requires a
	// logged-in user; the flags only add requirements on top
	if strings.Contains(strings.ToLower(ref), "authboss.middleware") {
		return Auth
	}
	if i := strings.IndexByte(ref, '='); i >= 0 {
		name, origin := ref[:i], ref[i+1:]
		if k, ok := c.override(name); ok {
			return k
		}
		if k := c.Classify(origin); k == Auth || k == NotAuth {
			return k
		}
		return c.Classify(name)
	}
	return classifyWords(words(ref))
}

// Optional reports a name that marks auth as optional (OptionalAuth,
// auth.Maybe), which a 401 somewhere in its body doesn't change.
func (c *Classifier) Optional(ref string) bool {
	if k, ok := c.override(ref); ok {
		return k != Auth
	}
	for _, w := range words(ref) {
		if optionalWords[w] {
			return true
		}
	}
	return false
}

// Overridden reports whether -auth or -not-auth names this reference.
func (c *Classifier) Overridden(ref string) bool {
	if _, ok := c.override(ref); ok {
		return true
	}
	if i := strings.IndexByte(ref, '='); i >= 0 {
		_, ok := c.override(ref[:i])
		return ok
	}
	return false
}

func (c *Classifier) override(ref string) (Kind, bool) {
	low := strings.ToLower(ref)
	base := low
	if i := strings.LastIndexByte(base, '.'); i >= 0 {
		base = base[i+1:]
	}
	for _, n := range c.notAuth {
		if n == low || n == base {
			return NotAuth, true
		}
	}
	for _, n := range c.auth {
		if n == low || n == base {
			return Auth, true
		}
	}
	return Unknown, false
}

// Vocabulary. Each set is whole lower-case words.
var (
	// Words that mean an auth check on their own.
	strongWords = set("auth", "authn", "authz", "authenticate", "authenticated", "authentication",
		"authenticator", "authorize", "authorized", "authorization", "authorizer", "authorizes",
		"jwt", "jwtauth", "jwtware", "echojwt", "bearer", "rbac", "abac", "acl", "casbin",
		"basicauth", "keyauth", "apikey", "loginrequired", "authrequired", "authmw", "authware",
		"protected")
	// Things an auth check is about.
	subjectWords = set("user", "users", "login", "logged", "signed", "signin", "session", "sessions",
		"token", "tokens", "role", "roles", "admin", "admins", "permission", "permissions", "perm",
		"perms", "scope", "scopes", "claim", "claims", "principal", "identity", "account",
		"owner", "member", "membership", "credential", "credentials", "policy", "policies",
		"oauth", "oauth2", "oidc", "openid", "saml", "sso", "privilege", "privileges", "superuser",
		"staff", "operator", "api_key", "key", "keys", "access", "group", "groups", "tenant",
		"client", "clients", "application", "app", "elevated")
	// Verbs that turn a subject into a check.
	actionWords = set("require", "required", "requires", "ensure", "ensures", "must", "check", "checks",
		"checker", "verify", "verifies", "verifier", "verified", "validate", "validates", "validator",
		"validated", "enforce", "enforcer", "enforces", "guard", "guarded", "protect", "restrict",
		"restricted", "only", "need", "needs", "is", "has", "can", "allow", "allowed", "gate", "gatekeeper",
		"deny", "control")
	// Hints that, alone, aren't enough to trust a name as a check.
	hintWords = set("session", "sessions", "token", "tokens", "oauth", "oauth2", "oidc", "openid",
		"saml", "sso", "guard", "enforce", "enforcer", "permission", "permissions", "policy", "rbac",
		"admin", "role", "roles", "acl", "scope", "scopes", "claims", "principal", "identity", "gatekeeper")
	// Modifiers that make an auth middleware non-enforcing: OptionalAuth.
	optionalWords = set("optional", "maybe", "try", "soft", "lenient", "lax", "anonymous", "anon",
		"skip", "bypass", "without", "noauth", "unauthenticated", "public", "allowanonymous",
		"ifpresent", "passthrough", "fake", "mock", "noop", "disabled", "disable", "none", "guest")
	// Verbs of middleware that reads auth state without rejecting anything
	// (session.Load, ParseToken), unless an enforcing verb is also present.
	loaderWords = set("load", "loader", "loads", "populate", "inject", "parse", "parser", "extract",
		"attach", "propagate", "forward", "refresh", "logout", "store")
	enforcingWords = set("require", "required", "requires", "ensure", "ensures", "must", "check", "checks",
		"verify", "verifies", "verifier", "validate", "validates", "validator", "enforce", "enforcer",
		"guard", "protect", "protected", "restrict", "restricted", "only", "deny", "gate")
	// Words that say nothing about what a middleware does.
	genericWords = set("middleware", "middlewares", "mw", "handler", "handle", "handlers", "func", "fn",
		"new", "config", "wrapper", "wrap", "filter", "interceptor", "plugin", "http", "h", "default",
		"v1", "v2", "v3", "v4", "v5", "with", "make", "build", "create", "get")
	// Words that mark middleware as something other than auth.
	notAuthWords = set("logger", "log", "logging", "recover", "recovery", "recoverer", "cors", "gzip",
		"compress", "compression", "deflate", "brotli", "requestid", "reqid", "ratelimit", "ratelimiter",
		"limiter", "limit", "throttle", "throttler", "bucket", "quota", "timeout", "bodylimit", "sizelimit",
		"static", "fileserver", "metrics", "prometheus", "tracing", "tracer", "trace", "otel", "telemetry",
		"cache", "nocache", "etag", "realip", "heartbeat", "cleanpath", "redirectslashes", "stripslashes",
		"csrf", "xsrf", "csp", "hsts", "helmet", "secure", "nosurf", "instrument", "instrumentation",
		"author", "authors", "authored", "authoring", "authorship", "profiler", "pprof", "debug",
		"i18n", "locale", "language", "theme", "decompress", "minify", "sentry", "newrelic", "datadog",
		"feature", "flag", "flags", "maintenance", "version", "nosniff", "accesslog")
)

// classifyWords applies the rules, most specific first:
//  1. a not-auth head word (Logger, RateLimit, tokenBucket), or a not-auth
//     word with no strong auth word (csrf.Protect) -> NotAuth
//  2. an optional modifier (OptionalAuth), or a loader verb with no
//     enforcing verb (session.Load), on auth words -> Weak
//  3. a strong word, or a check verb with an auth subject (RequireRole) -> Auth
//  4. an auth-ish hint alone (session, featureGuard) -> Weak
func classifyWords(ws []string) Kind {
	if len(ws) == 0 {
		return Unknown
	}
	// The head noun comes last in Go names: authLogger is a logger,
	// metricsBasicAuth is basic auth.
	head := ""
	for i := len(ws) - 1; i >= 0; i-- {
		if !genericWords[ws[i]] {
			head = ws[i]
			break
		}
	}
	if notAuthWords[head] {
		return NotAuth
	}
	strong, action, subject, hint, optional, loader, enforcing, other := false, false, false, false, false, false, false, false
	for i, w := range ws {
		other = other || notAuthWords[w]
		strong = strong || strongWords[w]
		action = action || actionWords[w]
		subject = subject || subjectWords[w]
		hint = hint || hintWords[w]
		optional = optional || optionalWords[w]
		loader = loader || loaderWords[w]
		enforcing = enforcing || enforcingWords[w]
		// "logged in", "signed in"; a bare "logged" middleware is the same check
		if (w == "logged" || w == "signed") && (i+1 < len(ws) && ws[i+1] == "in" || len(ws) == 1) {
			subject, action, enforcing = true, true, true
		}
		if w == "loggedin" || w == "signedin" {
			subject, action, enforcing = true, true, true
		}
		// "no auth", "skip auth"
		if w == "no" && i+1 < len(ws) && strongWords[ws[i+1]] {
			optional = true
		}
	}
	authy := strong || subject || hint
	if other && !strong {
		return NotAuth // quota.Enforce, csrf.Protect, featureGuard
	}
	switch {
	case optional && authy:
		return Weak
	case loader && !enforcing && authy:
		return Weak
	case strong:
		return Auth
	case action && subject:
		return Auth
	case hint:
		return Weak
	}
	return Unknown
}

// words splits an identifier reference into lower-case words.
// "middleware.RequireAuth" -> [middleware require auth];
// "JWTAuth" -> [jwt auth]; "OAuth2Middleware" -> [oauth2 middleware].
func words(ref string) []string {
	ref = strings.ReplaceAll(ref, "OAuth", "Oauth")
	ref = strings.ReplaceAll(ref, "OIDC", "Oidc")
	ref = strings.ReplaceAll(ref, "OpenID", "Openid")
	ref = strings.ReplaceAll(ref, "APIKey", "Apikey")
	ref = strings.ReplaceAll(ref, "ApiKey", "Apikey")
	var out []string
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			out = append(out, strings.ToLower(string(cur)))
			cur = cur[:0]
		}
	}
	rs := []rune(ref)
	for i, r := range rs {
		switch {
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			flush()
		case unicode.IsUpper(r):
			// boundary before an upper that follows a lower, or before the
			// last upper of an acronym followed by a lower (JWTAuth -> JWT Auth)
			if len(cur) > 0 {
				prev := rs[i-1]
				nextLower := i+1 < len(rs) && unicode.IsLower(rs[i+1])
				if unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && nextLower) {
					flush()
				}
			}
			cur = append(cur, r)
		case unicode.IsDigit(r):
			cur = append(cur, r) // oauth2, v1
		default:
			cur = append(cur, r)
		}
	}
	flush()
	return out
}

func set(xs ...string) map[string]bool {
	m := make(map[string]bool, len(xs))
	for _, x := range xs {
		m[x] = true
	}
	return m
}

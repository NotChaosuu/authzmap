package main

// allowlistApps: routes that are public by design (login, health, static,
// docs) and look-alikes that aren't.
var allowlistApps = map[string]appCase{
	"a01_allow_methods":    {M: []string{"DELETE /api/public/users/{id}", "DELETE /docs/{id}", "DELETE /health", "GET /ADMIN/export", "GET /Admin/users", "GET /admin/", "GET /api/%61dmin/users", "GET /api/public/../admin/users", "GET /api/users/{id}/token", "GET /debug/pprof/*", "GET /login/../admin/users", "GET /static/../admin/dump", "PATCH /api/public/users/{id}", "POST /api/v1/admin/users/register", "POST /api/v1/login/admin/promote", "POST /api/v1/users/{id}/password/reset", "POST /oauth/token/revoke-all", "POST /static/{file}", "POST /webhooks/github/../../admin/x", "PUT /static/{file}"}, Q: []string{"DELETE /api/notes/{id}", "GET /api/me", "POST /api/notes", "POST /health/reset", "PUT /api/public/settings", "PUT /api/settings"}, P: []string{"GET /api/public/stats", "GET /docs/index.html", "GET /health", "GET /metrics", "GET /static/app.js", "POST /login"}, exit: 1, checkExit: true},
	"a02_allow_nethttp122": {M: []string{"* /admin/", "DELETE /api/notes/{id}/{$}", "DELETE /assets/{path...}", "POST /admin/{path...}", "POST /static/{file}", "PUT /assets/{name}"}, Q: []string{"DELETE /api/notes/{id}", "GET /api/me", "POST /api/notes"}, P: []string{"* /static/", "GET /static/{path...}"}, exit: 1, checkExit: true},
	// the allowlist is method-aware and exact
	"d1_allow": {
		M: []string{"POST /api/projects/{id}/public", "DELETE /api/public/users/{id}", "DELETE /assets/{id}", "DELETE /webhooks/{id}",
			"POST /hooks/{id}/delete", "POST /api/jobs/{id}/callback", "POST /api/v1/users/{id}/files/{fileHash}",
			"DELETE /api/repos/{code}", "DELETE /api/countries/{countryCode}", "DELETE /oauth/{appID}", "POST /sso/{id}/delete",
			"DELETE /dist/{id}", "DELETE /api/users/{id}/openapi.json", "DELETE /{user}/repos/{id}",
			"DELETE /api/users/{id}/sessions/{sid}", "DELETE /docs/page", "DELETE /forgot-everything", "DELETE /api/sessions"},
		R: []string{"POST /static/upload", "POST /password/change", "POST /api/nodes/register", "POST /oauth/apps",
			"POST /api/swagger/config", "POST /swagger/config", "POST /.well-known/acme-reset",
			"POST /reset-password/admin-override", "POST /settings/token", "GET /api/users/export/callback"},
		P: []string{"GET /login", "GET /api/info", "GET /api/status", "POST /api/token", "POST /token", "POST /api/users/register"},
	},
	"d2_admin_entry": {
		M: []string{"POST /admin/impersonate/login", "POST /admin/jobs/callback", "POST /api/admin/sync/callback",
			"POST /internal/users/signin", "POST /admin/assets/*", "GET /api/admin/status"},
		R: []string{"POST /api/devices/register"}, P: []string{"GET /admin/static/*", "POST /api/v1/register"},
	},
	"q70_allowlist": {M: []string{"DELETE /health", "DELETE /oauth/{provider}", "DELETE /register", "DELETE /webhooks/{id}", "GET /api/v1/auth/token", "GET /status/../admin", "POST /admin/login"}, Q: []string{"GET /api/login/../users", "GET /api/me", "GET /version/{id}/download", "POST /api/public/lists", "POST /robots.txt", "POST /static/upload"}, P: []string{"GET /.well-known/openid-configuration", "GET /api/info", "GET /api/system/info", "GET /health", "GET /metrics", "GET /public/{id}", "GET /robots.txt", "GET /static/{path}", "GET /swagger/{path}", "POST /login", "POST /oauth/{provider}/callback", "POST /password/reset", "POST /register", "POST /webhooks/github"}, exit: 1, checkExit: true},
}

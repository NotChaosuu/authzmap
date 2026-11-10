# Changelog

## v1.0.0 (2026-11-09)

First release.

- `scan` reports HTTP routes reachable without an authorization check, ranked high / medium / low / info; `map` lists every route with the middleware in front of it.
- Frameworks: chi, gin, echo, fiber, gorilla/mux, hertz and net/http (Go 1.22 patterns included). grpc-gateway, gqlgen, connect-rpc and generated OpenAPI servers are found as one mounted handler.
- Middleware declared in the scanned code is read path by path, and its code wins over its name: auth that never rejects, rejects only sometimes, or is too branchy to verify is reported as such.
- Locals reassigned inside a branch (`ok = true` on a debug header, `tok = "k"` on a legacy flag) are followed per path, so a bypass through them shows as conditional auth.
- Handler builders that take the handler and its middleware (`chain.ToHandlerFunc(h, mw...)`, `Wrap(h, Logger, RequireAuth)`) credit that middleware.
- `map` has a status per route (`auth`, `public`, `weak`, `UNVERIFIED`, `COND-AUTH`, `NO-REJECT`, `NO-AUTH`) with a legend line, and the same `status` in JSON.
- `scan -format json` and `map -format json` share one header: `tool`, `version`, `mode`, `complete`, `files`, `partial_routes`, `problems`, `untraced`, `unseen`, `notes`. Lists are always arrays.
- Paths are printed with forward slashes on every OS.
- Scans that couldn't read something say so: unparsable files and routers it doesn't model make the scan incomplete (exit 3 with `-exit-code`), and registrations it couldn't trace mark their routes partial.

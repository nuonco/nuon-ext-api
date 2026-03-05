# Agents Plan

This repository holds a golang Nuon Extension.

Docs:
https://raw.githubusercontent.com/nuonco/nuon/refs/heads/fd/cli/implement-extensions/docs/guides/cli-extensions.mdx

This extension will be compiled and distributed as a binary via the releases.

1. read the docs.
2. initialize a go project called nuon-ext-api.

## Scope

We are concerned with the API client right now. The API TUI will be a future project which is only stubbed out at the
moment.

## Structure

This project is two things:

### API

an api client that can be used like this `nuon api path` and `nuon api path "payload"`. The paths are determined via the
swagger spec for the nuon public api (api.nuon.co).

This API has simple command and arg -based invocation but if a selection is required, it uses charm's bubble tea to
render a contextual TUI for selection.

### TUI

An API client TUI similar to https://github.com/darrenburns/posting but built for a specific version of the Nuon public
API.

The TUI will not be implemented at this time but its subcommand should be stubbed out.

The TUI is built with the latest charm bubbletea library.

---

## Approach Decision

### Options Considered

#### Spec-Driven Raw HTTP Client

Embed the Swagger 2.0 spec (`doc.json`) into the binary at compile time. Parse it at startup to build a route table
mapping API paths to HTTP methods, parameters, and schemas. The extension makes raw HTTP requests with auth from
environment variables and outputs JSON.

- The UX model (`nuon api <path> [payload]`) is inherently REST — raw HTTP maps directly
- 100% API coverage with zero per-endpoint code (~160+ endpoints covered automatically)
- Spec updates are trivial: re-embed the JSON, rebuild
- Minimal dependencies: stdlib `net/http` + `encoding/json` + bubbletea for interactive selection
- No go-swagger/go-openapi dependency tree (the SDK pulls in ~50 transitive deps)
- Small binary size
- Can be built incrementally

---

## Architecture

### Extension Protocol

The Nuon CLI invokes extensions as subprocesses. For this extension (`api`), the user runs:

```
nuon api /v1/apps
nuon api /v1/apps '{"name":"my-app"}'
```

The CLI calls the binary with raw args and these environment variables:

| Variable           | Description                                |
| ------------------ | ------------------------------------------ |
| `NUON_API_URL`     | API endpoint (e.g., `https://api.nuon.co`) |
| `NUON_API_TOKEN`   | Bearer token for authentication            |
| `NUON_ORG_ID`      | Current organization ID                    |
| `NUON_APP_ID`      | Current app ID (if set)                    |
| `NUON_INSTALL_ID`  | Current install ID (if set)                |
| `NUON_CONFIG_FILE` | Path to CLI config file                    |
| `NUON_EXT_NAME`    | Extension name (`api`)                     |
| `NUON_EXT_DIR`     | Extension directory path                   |

The binary receives stdin/stdout/stderr directly from the parent process, enabling interactive TUI.

### Project Layout

```
nuon-ext-api/
├── nuon-ext.toml              # Extension manifest
├── main.go                    # Entry point: parse env, dispatch subcommands
├── cmd/
│   ├── root.go                # Root cobra command
│   ├── api.go                 # `api` subcommand (or default when invoked as extension)
│   └── tui.go                 # `tui` subcommand stub
├── internal/
│   ├── config/
│   │   └── config.go          # Read env vars into a config struct
│   ├── spec/
│   │   ├── spec.go            # Parse embedded swagger spec, build route table
│   │   └── route.go           # Route type: path, method, params, request/response schemas
│   ├── client/
│   │   └── client.go          # HTTP client: auth headers, request execution, response handling
│   ├── dispatch/
│   │   └── dispatch.go        # Match user input path to spec routes, resolve method
│   ├── selector/
│   │   ├── selector.go        # Bubbletea model for interactive selection
│   │   └── resolve.go         # Resolve path params ({app_id} etc.) via API + selector
│   └── output/
│       └── output.go          # JSON pretty-printing, optional table/compact modes
├── spec/
│   └── doc.json               # Embedded swagger spec (go:embed)
├── go.mod
├── go.sum
├── .goreleaser.yml            # GoReleaser config for cross-platform binaries
└── .github/
    └── workflows/
        └── release.yml        # GitHub Actions: build + release on tag
```

### Core Flow

```
User runs: nuon api /v1/apps/{app_id}

1. main.go reads env vars → config.Config{APIURL, Token, OrgID, ...}
2. cmd/api.go receives args: ["/v1/apps/{app_id}"]
3. dispatch.Resolve(path, payload) →
   a. Looks up "/v1/apps/{app_id}" in the spec route table
   b. No payload → selects GET method
   c. Detects {app_id} path parameter
4. selector.ResolvePathParams(route, knownContext) →
   a. Checks if NUON_APP_ID is set → uses it for {app_id}
   b. If not, calls GET /v1/apps to fetch list
   c. Renders bubbletea list selector → user picks an app
   d. Returns resolved path: /v1/apps/app-abc123
5. client.Do(method, resolvedPath, payload, headers) →
   a. Builds request with Authorization: Bearer <token>, X-Nuon-Org-ID: <org>
   b. Sends HTTP request
   c. Returns response body + status
6. output.Print(responseJSON) → pretty-printed to stdout
```

### Method Inference

When the user provides a path:

| Condition                                           | Inferred Method                         |
| --------------------------------------------------- | --------------------------------------- |
| No payload, path has GET in spec                    | GET                                     |
| Payload provided, path has POST in spec             | POST                                    |
| Payload provided, path has PATCH but no POST        | PATCH                                   |
| Payload provided, path has PUT but no POST or PATCH | PUT                                     |
| Ambiguous (multiple write methods)                  | Error with hint to use `-X METHOD` flag |

The user can always override with `-X METHOD` (e.g., `nuon api -X DELETE /v1/apps/{app_id}`).

### Path Parameter Resolution

Path parameters like `{app_id}`, `{install_id}`, `{component_id}` are resolved in priority order:

1. **Literal value in path**: `nuon api /v1/apps/app-abc123` → no resolution needed
2. **Environment context**: `{app_id}` → `NUON_APP_ID`, `{install_id}` → `NUON_INSTALL_ID`, `{org_id}` → `NUON_ORG_ID`
3. **Interactive selection**: Fetch the parent resource list and render a bubbletea selector

For interactive selection, the spec tells us the parameter type, and we know the list endpoints:

- `{app_id}` → `GET /v1/apps`
- `{install_id}` → `GET /v1/installs` (or scoped to app if app_id known)
- `{component_id}` → `GET /v1/components` (or scoped to app)
- Other params → prompt for text input

### Swagger Spec Embedding

The spec is embedded at compile time:

```go
//go:embed spec/doc.json
var specJSON []byte
```

The spec file is copied from `https://api.nuon.co/docs/doc.json` (or locally from
`~/nuon/nuon/services/ctl-api/docs/public/swagger.json`) during development. A Makefile target handles this.

At startup, parse the spec into a route table:

```go
type Route struct {
    Path       string            // e.g., "/v1/apps/{app_id}"
    Method     string            // e.g., "GET"
    OperationID string           // e.g., "GetApp"
    Summary    string            // Human-readable description
    PathParams []Param           // Parameters in the path
    QueryParams []Param          // Query string parameters
    BodySchema *Schema           // Request body schema (for POST/PUT/PATCH)
}
```

### Dependencies

Minimal dependency set:

| Dependency                           | Purpose                                  |
| ------------------------------------ | ---------------------------------------- |
| `github.com/spf13/cobra`             | CLI framework (consistent with nuon CLI) |
| `github.com/charmbracelet/bubbletea` | Interactive TUI selection                |
| `github.com/charmbracelet/lipgloss`  | TUI styling                              |
| `github.com/charmbracelet/bubbles`   | TUI components (list, textinput)         |
| stdlib `net/http`                    | HTTP client                              |
| stdlib `encoding/json`               | JSON parse/output                        |

No go-swagger. No go-openapi. No generated code.

### Spec Parsing

The Swagger 2.0 spec is a well-known JSON format. We only need to parse:

- `paths` → iterate all paths and methods to build routes
- `paths.{path}.{method}.parameters` → extract path/query/body params
- `paths.{path}.{method}.operationId` → operation name for display
- `paths.{path}.{method}.summary` → description for help text
- `definitions` → only needed if we want to show schema info for `--help`

This is ~200 lines of parsing code using `encoding/json` and a few struct types matching the Swagger 2.0 schema.

---

## Implementation Phases

### Phase 1: Project Scaffold

- [ ] Initialize Go module (`nuon-ext-api`)
- [ ] Create `nuon-ext.toml` manifest
- [ ] Create `main.go` with cobra root command
- [ ] Create `cmd/api.go` and `cmd/tui.go` (stub)
- [ ] Create `internal/config/config.go` (read env vars)
- [ ] Copy and embed `spec/doc.json`
- [ ] Add `.goreleaser.yml` for cross-platform builds
- [ ] Add `.gitignore` with Go-specific rules

### Phase 2: Spec Parser + Route Table

- [ ] Implement `internal/spec/spec.go` — parse embedded swagger JSON
- [ ] Implement `internal/spec/route.go` — Route type with path, method, params
- [ ] Build route table from spec at startup
- [ ] Add `nuon api --list` — interactive bubbletea endpoint browser with fuzzy search

### Phase 3: HTTP Client + Basic Requests

- [ ] Implement `internal/client/client.go` — HTTP client with auth headers
- [ ] Implement `internal/dispatch/dispatch.go` — match user path to route, infer method
- [ ] Implement `internal/output/output.go` — pretty-print JSON
- [ ] End-to-end: `nuon api /v1/apps` returns pretty JSON

### Phase 4: Path Parameter Resolution + Interactive Selection

- [ ] Implement `internal/selector/resolve.go` — resolve params from env or interactively
- [ ] Implement `internal/selector/selector.go` — bubbletea list model
- [ ] Wire up: `nuon api /v1/apps/{app_id}` resolves via env or selector
- [ ] Support contextual scoping (e.g., use `NUON_APP_ID` when available)

### Phase 5: Polish + Release

- [ ] Add `-X METHOD` flag for explicit method override
- [ ] Add `-q` / `--query` for query parameters
- [ ] Add `--raw` flag for unformatted JSON output
- [ ] Add `--help` per-endpoint (show params, schema info from spec)
- [ ] GitHub Actions release workflow
- [ ] Test with `nuon ext install ./nuon-ext-api` locally

---

## Design Decisions

1. **Endpoint listing UX**: `nuon api` with no args shows help. `nuon api --list` launches an interactive bubbletea
   endpoint browser with fuzzy search.
2. **Output format**: JSON only. This is an API client. The TUI (future) may support additional output formats.
3. **Spec freshness**: The spec is set at compile time via `go:embed`. At release time, the GitHub Action checks
   `api.nuon.co/version` and uses that version string to tag the release, keeping the extension version in sync with
   the API version it targets.

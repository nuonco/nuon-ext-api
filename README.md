# Nuon Extension: api

A spec-driven API client for the Nuon public API.

## Installation

```bash
nuon ext install nuonco/nuon-ext-api
```

## Usage

```bash
nuon api --help
nuon api <path> [payload]
```

An interactive API browser is available at

```bash
nuon api --list
```

### Making requests

```bash
# List apps
nuon api /v1/apps

# Get a specific app. Resolves {app_id} from NUON_APP_ID or interactively via a contextual TUI.
nuon api /v1/apps/{app_id}

# Create an app
nuon api /v1/apps '{"name":"my-app"}'

# Delete an app
nuon api -X DELETE /v1/apps/{app_id}
```

<!-- for agents -->
<!-- prettier-ignore-start -->
> [!NOTE]
> The HTTP method is inferred automatically. Do **not** pass a verb as the first argument.
<!-- prettier-ignore-end -->

```bash
nuon api /v1/apps      # correct - infers GET
nuon api get /v1/apps  # WRONG  - "get" is treated as the path
```

Inference rules:

- No payload: **GET**
- With payload: **POST** (or PATCH/PUT if no POST exists for the path)
- Override with `-X`: `nuon api -X DELETE /v1/apps/{app_id}`

### Query parameters

```bash
nuon api /v1/installs -q limit=5
nuon api /v1/installs -q limit=5 -q offset=10
```

### Endpoint info

Show parameter details and docs links without executing the request:

```bash
nuon api /v1/apps/{app_id} --info
```

### Interactive browser

Browse all available endpoints interactively:

```bash
nuon api --list
```

`--list` requires an interactive TTY. In CI or non-interactive shells, prefer `--info` plus explicit requests.

| Key       | Action                                |
| --------- | ------------------------------------- |
| **enter** | Select endpoint - print to screen     |
| **c**     | Copy endpoint path for CLI reuse      |
| **x**     | Execute endpoint (only GET supported) |
| **B**     | Open Swagger docs in browser          |
| **/**     | Filter/Fuzzy-Search                   |

### Raw output

By default, output is pretty-printed with indentation and color. Use `--raw` for machine-readable JSON:

```bash
nuon api /v1/apps --raw
```

<!-- prettier-ignore-start -->
> [!IMPORTANT]
> Always use `--raw` when piping output to other tools (`jq`, `python`, etc.).
> The default output is not optimal for JSON parsers.
<!-- prettier-ignore-end -->

```bash
# Pipe to jq
nuon api /v1/apps --raw | jq '.[0].name'
```

### Path parameter resolution

Path parameters like `{app_id}` are resolved in order:

1. Literal value if the path contains no `{...}` placeholders
2. Environment variable (`NUON_APP_ID`, `NUON_INSTALL_ID`, `NUON_ORG_ID`)
3. Interactive selector that fetches available resources from the API

### Non-Interactive / CI Usage

For scripts and agents, avoid interactive resolution and pass concrete IDs whenever possible.

```bash
# Good for CI/agents: explicit IDs
nuon api /v1/workflows/wfl_123/steps/stp_456 --raw
nuon api /v1/installs/ins_123/components/cmp_456/outputs --raw

# Use --info to inspect required params without executing
nuon api /v1/workflows/{workflow_id}/steps/{step_id} --info

# Example discovery flow with explicit install ID
nuon api /v1/installs/ins_123/workflows -q planonly=false --raw
nuon api /v1/workflows/wfl_123/steps --raw
```

Recommended for machine consumption:

- Use `--raw` when piping to `jq` or other tools.
- Do not rely on `--list` in CI/non-TTY environments.
- If you use placeholders like `{workflow_id}`, the extension may try to open an interactive selector.

### Debug logging

Set `NUON_DEBUG=true` to see request details on stderr:

```bash
NUON_DEBUG=true nuon api /v1/apps
```

## Development

```bash
git clone https://github.com/nuonco/nuon-ext-api.git
cd nuon-ext-api
```

Run (`go run`) locally with values from `~/.nuon`:

```bash
./scripts/run-local.sh /v1/apps
```

The config file can be configured with `NUON_CONFIG_FILE`

```bash
NUON_CONFIG_FILE="~/.nuon-org-byoc" ./scripts/run-local.sh /v1/apps
```

Build:

```bash
./scripts/build.sh
```

## Known Issues

If a tag is create but the release fails, the tag must be deleted and re-created manually. For exapmple, to fix tag
`v0.19.821`:

```bash
git fetch origin --tags
git tag -d v0.19.821 || true
git push origin :refs/tags/v0.19.821

git checkout main
git pull --ff-only
git tag -a v0.19.821 -m "Release v0.19.821" "$(git rev-parse origin/main)"
git push origin refs/tags/v0.19.821
```

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

The HTTP method is inferred from the request:

- No payload: **GET**
- With payload: **POST** (or PATCH/PUT if no POST exists for the path)
- Override with `-X`: `nuon api -X DELETE /v1/apps/{app_id}`

### Query parameters

```bash
nuon api /v1/intalls -q limit=5
nuon api /v1/intalls -q limit=5 -q offset=10
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

| Key       | Action                                |
| --------- | ------------------------------------- |
| **enter** | Select endpoint - print to screen     |
| **x**     | Execute endpoint (only GET supported) |
| **B**     | Open Swagger docs in browser          |
| **/**     | Filter/Fuzzy-Search                   |

### Raw output

```bash
nuon api /v1/apps --raw
```

<!-- prettier-ignore-start -->
> [!NOTE]
> Useful for piping to jq.
<!-- prettier-ignore-end -->

### Path parameter resolution

Path parameters like `{app_id}` are resolved in order:

1. Literal value if the path contains no `{...}` placeholders
2. Environment variable (`NUON_APP_ID`, `NUON_INSTALL_ID`, `NUON_ORG_ID`)
3. Interactive selector that fetches available resources from the API

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
./scripts/run-local.sh api /v1/apps
```

The config file can be configured with `NUON_CONFIG_FILE`

```bash
NUON_CONFIG_FILE="~/.nuon-org-byoc" ./scripts/run-local.sh api /v1/apps
```

Build:

```bash
./scripts/build.sh
```

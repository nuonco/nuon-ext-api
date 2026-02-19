package resolve

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nuonco/nuon-ext-api/internal/client"
	"github.com/nuonco/nuon-ext-api/internal/config"
	"github.com/nuonco/nuon-ext-api/internal/pkg/tui/selector"
)

// envMap maps path parameter names to environment variable values.
var envMap = map[string]func(cfg *config.Config) string{
	"app_id":     func(cfg *config.Config) string { return cfg.AppID },
	"install_id": func(cfg *config.Config) string { return cfg.InstallID },
	"org_id":     func(cfg *config.Config) string { return cfg.OrgID },
}

// listEndpoints maps path parameter names to the API endpoint that lists resources of that type.
// The value is the path template — if it contains {app_id}, the app_id must be resolved first.
var listEndpoints = map[string]string{
	"app_id":                  "/v1/apps",
	"install_id":              "/v1/installs",
	"component_id":            "/v1/components",
	"org_id":                  "/v1/orgs",
	"action_workflow_id":      "/v1/action-workflows",
	"workflow_id":             "/v1/workflows",
	"vcs_connection_id":       "/v1/vcs-connections",
}

// PathParams resolves all {param} placeholders in a path.
// Priority: literal values already in the path > env vars > interactive selection.
// Returns the fully resolved path.
func PathParams(path string, cfg *config.Config, c *client.Client) (string, error) {
	if !strings.Contains(path, "{") {
		return path, nil
	}

	// Track resolved values so scoped lookups can use them (e.g., app_id for install listing)
	resolved := make(map[string]string)

	parts := strings.Split(path, "/")
	for i, part := range parts {
		if !strings.HasPrefix(part, "{") || !strings.HasSuffix(part, "}") {
			continue
		}

		paramName := part[1 : len(part)-1]

		// 1. Check env var
		if fn, ok := envMap[paramName]; ok {
			if val := fn(cfg); val != "" {
				parts[i] = val
				resolved[paramName] = val
				continue
			}
		}

		// 2. Interactive selection
		val, err := selectParam(paramName, cfg, c, resolved)
		if err != nil {
			return "", err
		}
		parts[i] = val
		resolved[paramName] = val
	}

	return strings.Join(parts, "/"), nil
}

func selectParam(paramName string, cfg *config.Config, c *client.Client, resolved map[string]string) (string, error) {
	listPath, ok := listEndpoints[paramName]
	if !ok {
		return "", fmt.Errorf("cannot resolve {%s}: no list endpoint known and no env var set", paramName)
	}

	resp, err := c.Do("GET", listPath, "")
	if err != nil {
		return "", fmt.Errorf("fetching resources for {%s}: %w", paramName, err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("fetching resources for {%s}: HTTP %d", paramName, resp.StatusCode)
	}

	resources, err := parseResources(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parsing resources for {%s}: %w", paramName, err)
	}

	result, err := selector.Run(paramName, resources)
	if err != nil {
		return "", err
	}
	if !result.Selected {
		return "", fmt.Errorf("no selection made for {%s}", paramName)
	}

	return result.ID, nil
}

// parseResources extracts id+name from a JSON array of objects.
func parseResources(data []byte) ([]selector.Resource, error) {
	var items []map[string]any
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	resources := make([]selector.Resource, 0, len(items))
	for _, item := range items {
		id, _ := item["id"].(string)
		if id == "" {
			continue
		}

		// Use the best available display name
		name := stringField(item, "display_name", "name", "id")

		resources = append(resources, selector.Resource{
			ID:   id,
			Name: name,
		})
	}

	return resources, nil
}

func stringField(obj map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := obj[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

package spec

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	embeddedSpec "github.com/nuonco/nuon-ext-api/spec"
)

// API holds the parsed route table from the swagger spec.
type API struct {
	Version string
	Routes  []Route
	byPath  map[string][]Route // path template → routes for all methods
}

// Parse reads the embedded swagger spec and builds the route table.
func Parse() (*API, error) {
	var raw swaggerDoc
	if err := json.Unmarshal(embeddedSpec.JSON, &raw); err != nil {
		return nil, fmt.Errorf("parsing swagger spec: %w", err)
	}

	api := &API{
		Version: raw.Info.Version,
		byPath:  make(map[string][]Route),
	}

	for path, methods := range raw.Paths {
		for method, op := range methods {
			method = strings.ToUpper(method)
			if !isHTTPMethod(method) {
				continue
			}

			route := Route{
				Path:        path,
				Method:      method,
				OperationID: op.OperationID,
				Summary:     op.Summary,
			}
			if len(op.Tags) > 0 {
				route.Tag = op.Tags[0]
			}

			for _, p := range op.Parameters {
				param := Param{
					Name:        p.Name,
					In:          p.In,
					Type:        p.Type,
					Required:    p.Required,
					Description: p.Description,
					Default:     p.Default,
				}
				switch p.In {
				case "path":
					route.PathParams = append(route.PathParams, param)
				case "query":
					route.QueryParams = append(route.QueryParams, param)
				case "body":
					route.HasBody = true
					if p.Schema != nil {
						route.BodySchema = p.Schema.Ref
					}
				}
			}

			api.Routes = append(api.Routes, route)
			api.byPath[path] = append(api.byPath[path], route)
		}
	}

	sort.Slice(api.Routes, func(i, j int) bool {
		if api.Routes[i].Path != api.Routes[j].Path {
			return api.Routes[i].Path < api.Routes[j].Path
		}
		return methodOrder(api.Routes[i].Method) < methodOrder(api.Routes[j].Method)
	})

	return api, nil
}

// Lookup finds all routes matching a given path (with or without concrete param values).
// It first tries an exact template match, then tries matching against route patterns.
func (a *API) Lookup(inputPath string) []Route {
	// Exact template match (e.g., user typed "/v1/apps/{app_id}")
	if routes, ok := a.byPath[inputPath]; ok {
		return routes
	}

	// Pattern match (e.g., user typed "/v1/apps/abc123" → matches "/v1/apps/{app_id}")
	var matches []Route
	for _, route := range a.Routes {
		if ok, _ := route.MatchesPath(inputPath); ok {
			matches = append(matches, route)
		}
	}

	// Deduplicate by method (multiple routes can match the same path)
	seen := make(map[string]bool)
	var unique []Route
	for _, r := range matches {
		key := r.Method + " " + r.Path
		if !seen[key] {
			seen[key] = true
			unique = append(unique, r)
		}
	}

	return unique
}

// LookupByMethod finds a specific route for a path and method.
func (a *API) LookupByMethod(inputPath, method string) *Route {
	method = strings.ToUpper(method)
	for _, r := range a.Lookup(inputPath) {
		if r.Method == method {
			return &r
		}
	}
	return nil
}

// swagger 2.0 JSON structures — only the fields we need

type swaggerDoc struct {
	Info  swaggerInfo                          `json:"info"`
	Paths map[string]map[string]swaggerOp      `json:"paths"`
}

type swaggerInfo struct {
	Version string `json:"version"`
}

type swaggerOp struct {
	OperationID string           `json:"operationId"`
	Summary     string           `json:"summary"`
	Tags        []string         `json:"tags"`
	Parameters  []swaggerParam   `json:"parameters"`
}

type swaggerParam struct {
	Name        string         `json:"name"`
	In          string         `json:"in"`
	Type        string         `json:"type"`
	Required    bool           `json:"required"`
	Description string         `json:"description"`
	Default     any            `json:"default"`
	Schema      *swaggerSchema `json:"schema"`
}

type swaggerSchema struct {
	Ref string `json:"$ref"`
}

func isHTTPMethod(m string) bool {
	switch m {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
		return true
	}
	return false
}

func methodOrder(m string) int {
	switch m {
	case "GET":
		return 0
	case "POST":
		return 1
	case "PUT":
		return 2
	case "PATCH":
		return 3
	case "DELETE":
		return 4
	}
	return 5
}

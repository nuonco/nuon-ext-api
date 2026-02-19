package dispatch

import (
	"fmt"
	"strings"

	"github.com/nuonco/nuon-ext-api/internal/client"
	"github.com/nuonco/nuon-ext-api/internal/config"
	"github.com/nuonco/nuon-ext-api/internal/resolve"
	"github.com/nuonco/nuon-ext-api/internal/spec"
)

// Request represents a resolved API request ready for execution.
type Request struct {
	Route   spec.Route
	Path    string // resolved path with concrete param values
	Method  string
	Payload string // raw JSON body (empty for GET/DELETE)
}

// Resolve takes user input (path, optional payload, optional method override)
// and matches it against the spec to produce an executable Request.
// If the path contains {param} placeholders, they are resolved via env vars
// or interactive selection.
func Resolve(api *spec.API, inputPath, payload, methodOverride string, cfg *config.Config, c *client.Client) (*Request, error) {
	// First, look up the route using the raw input (may contain {param} templates)
	routes := api.Lookup(inputPath)
	if len(routes) == 0 {
		return nil, fmt.Errorf("no endpoint found for path: %s", inputPath)
	}

	method := inferMethod(routes, payload, methodOverride)
	if method == "" {
		available := make([]string, len(routes))
		for i, r := range routes {
			available[i] = r.Method
		}
		return nil, fmt.Errorf(
			"ambiguous method for %s (available: %s) — use -X to specify",
			inputPath, strings.Join(available, ", "),
		)
	}

	var matched *spec.Route
	for _, r := range routes {
		if r.Method == method {
			matched = &r
			break
		}
	}
	if matched == nil {
		return nil, fmt.Errorf("method %s not available for path: %s", method, inputPath)
	}

	// Resolve path parameters.
	// If the input path still contains {param} placeholders, resolve them
	// via env vars or interactive selection. Otherwise use the input as-is.
	resolvedPath := inputPath
	if strings.Contains(inputPath, "{") {
		var err error
		resolvedPath, err = resolve.PathParams(matched.Path, cfg, c)
		if err != nil {
			return nil, err
		}
	}

	return &Request{
		Route:   *matched,
		Path:    resolvedPath,
		Method:  method,
		Payload: payload,
	}, nil
}

// inferMethod determines the HTTP method from context.
func inferMethod(routes []spec.Route, payload, override string) string {
	if override != "" {
		return strings.ToUpper(override)
	}

	hasPayload := payload != ""

	if !hasPayload {
		// No payload — prefer GET
		for _, r := range routes {
			if r.Method == "GET" {
				return "GET"
			}
		}
		// No GET available — if there's only one method, use it
		if len(routes) == 1 {
			return routes[0].Method
		}
		return ""
	}

	// Has payload — prefer POST, then PATCH, then PUT
	for _, prefer := range []string{"POST", "PATCH", "PUT"} {
		for _, r := range routes {
			if r.Method == prefer {
				return prefer
			}
		}
	}

	// Payload provided but no write method — if there's only one method, use it
	if len(routes) == 1 {
		return routes[0].Method
	}

	return ""
}

package dispatch

import (
	"fmt"
	"strings"

	"github.com/nuonco/nuon-ext-api/internal/spec"
)

// Request represents a resolved API request ready for execution.
type Request struct {
	Route      spec.Route
	Path       string // resolved path with concrete param values
	Method     string
	Payload    string            // raw JSON body (empty for GET/DELETE)
	PathParams map[string]string // extracted path parameter values
}

// Resolve takes user input (path, optional payload, optional method override)
// and matches it against the spec to produce an executable Request.
func Resolve(api *spec.API, inputPath, payload, methodOverride string) (*Request, error) {
	routes := api.Lookup(inputPath)
	if len(routes) == 0 {
		return nil, fmt.Errorf("no endpoint found for path: %s", inputPath)
	}

	// Extract path params from the match
	var pathParams map[string]string
	if len(routes) > 0 {
		_, pathParams = routes[0].MatchesPath(inputPath)
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

	// Resolve the actual path — if user passed a template, keep it for now
	// (Phase 4 will handle interactive resolution of {param} placeholders)
	resolvedPath := inputPath
	if pathParams == nil {
		pathParams = make(map[string]string)
	}

	return &Request{
		Route:      *matched,
		Path:       resolvedPath,
		Method:     method,
		Payload:    payload,
		PathParams: pathParams,
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

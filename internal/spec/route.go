package spec

import "strings"

// Route represents a single API endpoint (one method on one path).
type Route struct {
	Path        string  // e.g., "/v1/apps/{app_id}"
	Method      string  // e.g., "GET"
	OperationID string  // e.g., "GetApp"
	Summary     string  // Human-readable description
	Tag         string  // Primary tag (e.g., "apps")
	PathParams  []Param // Parameters in the path
	QueryParams []Param // Query string parameters
	HasBody     bool    // Whether the endpoint accepts a request body
	BodySchema  string  // $ref for the body schema (e.g., "#/definitions/service.CreateAppRequest")
}

// Param represents a single API parameter.
type Param struct {
	Name        string
	In          string // "path", "query", "body"
	Type        string // "string", "integer", etc.
	Required    bool
	Description string
	Default     any
}

// DisplayName returns a short display string like "GET /v1/apps".
func (r Route) DisplayName() string {
	return r.Method + " " + r.Path
}

// DocsURL returns the Swagger UI URL for this endpoint.
func (r Route) DocsURL(baseURL string) string {
	return strings.TrimRight(baseURL, "/") + "/docs/index.html#/" + r.Tag + "/" + r.OperationID
}

// MatchesPath checks if a user-provided path matches this route's template.
// For example, "/v1/apps/abc123" matches "/v1/apps/{app_id}".
// Returns true and the extracted path parameter values if it matches.
func (r Route) MatchesPath(inputPath string) (bool, map[string]string) {
	routeParts := strings.Split(strings.Trim(r.Path, "/"), "/")
	inputParts := strings.Split(strings.Trim(inputPath, "/"), "/")

	if len(routeParts) != len(inputParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i, rp := range routeParts {
		if strings.HasPrefix(rp, "{") && strings.HasSuffix(rp, "}") {
			paramName := rp[1 : len(rp)-1]
			params[paramName] = inputParts[i]
		} else if rp != inputParts[i] {
			return false, nil
		}
	}

	return true, params
}

// HasUnresolvedParams returns true if the path still contains {param} placeholders.
func (r Route) HasUnresolvedParams(path string) bool {
	return strings.Contains(path, "{") && strings.Contains(path, "}")
}

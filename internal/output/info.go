package output

import (
	"fmt"
	"strings"

	"github.com/nuonco/nuon-ext-api/internal/spec"
)

// PrintEndpointInfo displays detailed information about all routes matching a path.
func PrintEndpointInfo(routes []spec.Route, apiURL string) {
	for i, r := range routes {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s %s\n", r.Method, r.Path)
		if r.Summary != "" {
			fmt.Printf("  %s\n", r.Summary)
		}
		fmt.Printf("  Operation: %s\n", r.OperationID)
		fmt.Printf("  Docs:      %s\n", r.DocsURL(apiURL))

		if len(r.PathParams) > 0 {
			fmt.Println("  Path params:")
			for _, p := range r.PathParams {
				printParam(p)
			}
		}

		if len(r.QueryParams) > 0 {
			fmt.Println("  Query params:")
			for _, p := range r.QueryParams {
				printParam(p)
			}
		}

		if r.HasBody {
			schema := r.BodySchema
			if strings.HasPrefix(schema, "#/definitions/") {
				schema = strings.TrimPrefix(schema, "#/definitions/")
			}
			fmt.Printf("  Body:      %s\n", schema)
		}
	}
}

func printParam(p spec.Param) {
	req := ""
	if p.Required {
		req = " (required)"
	}
	def := ""
	if p.Default != nil {
		def = fmt.Sprintf(" [default: %v]", p.Default)
	}
	desc := ""
	if p.Description != "" {
		desc = " — " + p.Description
	}
	fmt.Printf("    %-20s %s%s%s%s\n", p.Name, p.Type, req, def, desc)
}

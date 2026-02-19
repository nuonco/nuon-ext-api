package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nuonco/nuon-ext-api/internal/pkg/tui/browser"
	"github.com/nuonco/nuon-ext-api/internal/spec"
)

var api *spec.API

func apiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api <path> [payload]",
		Short: "Make API requests to the Nuon public API",
		Long: `Make API requests to the Nuon public API.

The HTTP method is inferred from the request:
  - No payload: GET
  - With payload: POST (or PATCH/PUT if no POST exists for the path)

Override the method with -X:
  nuon api -X DELETE /v1/apps/{app_id}

Examples:
  nuon api /v1/apps
  nuon api /v1/apps '{"name":"my-app"}'
  nuon api /v1/apps/{app_id}
  nuon api --list`,
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: false,
		PersistentPreRunE:  initAPI,
		RunE:               runAPI,
	}

	cmd.Flags().StringP("method", "X", "", "HTTP method override (GET, POST, PUT, PATCH, DELETE)")
	cmd.Flags().Bool("list", false, "Browse available API endpoints interactively")
	cmd.Flags().Bool("raw", false, "Output raw JSON without formatting")

	return cmd
}

func initAPI(cmd *cobra.Command, args []string) error {
	var err error
	api, err = spec.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse API spec: %w", err)
	}
	return nil
}

func runAPI(cmd *cobra.Command, args []string) error {
	showList, _ := cmd.Flags().GetBool("list")
	if showList {
		result, err := browser.Run(api)
		if err != nil {
			return err
		}
		if result.Selected {
			fmt.Println(result.Route.DisplayName())
		}
		return nil
	}

	if len(args) == 0 {
		return cmd.Help()
	}

	// TODO: Phase 3+ — dispatch, HTTP client
	path := args[0]
	routes := api.Lookup(path)
	if len(routes) == 0 {
		return fmt.Errorf("no endpoint found for path: %s", path)
	}

	fmt.Printf("matched %d route(s) for %s:\n", len(routes), path)
	for _, r := range routes {
		fmt.Printf("  %s %s (%s) — %s\n", r.Method, r.Path, r.OperationID, r.Summary)
	}

	return nil
}

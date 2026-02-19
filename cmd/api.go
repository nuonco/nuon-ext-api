package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
		RunE:               runAPI,
	}

	cmd.Flags().StringP("method", "X", "", "HTTP method override (GET, POST, PUT, PATCH, DELETE)")
	cmd.Flags().Bool("list", false, "Browse available API endpoints interactively")
	cmd.Flags().Bool("raw", false, "Output raw JSON without formatting")

	return cmd
}

func runAPI(cmd *cobra.Command, args []string) error {
	list, _ := cmd.Flags().GetBool("list")
	if list {
		fmt.Println("interactive endpoint browser not yet implemented")
		return nil
	}

	if len(args) == 0 {
		return cmd.Help()
	}

	// TODO: Phase 2+ — spec parsing, dispatch, HTTP client
	fmt.Printf("path: %s\n", args[0])
	if len(args) > 1 {
		fmt.Printf("payload: %s\n", args[1])
	}

	return nil
}

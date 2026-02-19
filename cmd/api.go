package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nuonco/nuon-ext-api/internal/client"
	"github.com/nuonco/nuon-ext-api/internal/dispatch"
	"github.com/nuonco/nuon-ext-api/internal/output"
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

	raw, _ := cmd.Flags().GetBool("raw")
	methodOverride, _ := cmd.Flags().GetString("method")

	path := args[0]
	var payload string
	if len(args) > 1 {
		payload = args[1]
	}

	req, err := dispatch.Resolve(api, path, payload, methodOverride)
	if err != nil {
		return err
	}

	c := client.New(cfg)
	resp, err := c.Do(req.Method, req.Path, req.Payload)
	if err != nil {
		return err
	}

	return output.Print(resp, raw)
}

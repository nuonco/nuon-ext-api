package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nuonco/nuon-ext-api/internal/client"
	"github.com/nuonco/nuon-ext-api/internal/config"
	"github.com/nuonco/nuon-ext-api/internal/dispatch"
	"github.com/nuonco/nuon-ext-api/internal/output"
	"github.com/nuonco/nuon-ext-api/internal/pkg/tui/browser"
	"github.com/nuonco/nuon-ext-api/internal/spec"
)

var (
	cfg *config.Config
	api *spec.API
)

func Execute() {
	cfg = config.Load()

	root := &cobra.Command{
		Use:   "nuon-ext-api <path> [payload]",
		Short: "API client for the Nuon public API",
		Long: `Make API requests to the Nuon public API.

The HTTP method is inferred from the request:
  - No payload: GET
  - With payload: POST (or PATCH/PUT if no POST exists for the path)

Override the method with -X:
  nuon api -X DELETE /v1/apps/{app_id}

Examples:
  nuon api /v1/apps
  nuon api /v1/apps -q limit=5
  nuon api /v1/apps '{"name":"my-app"}'
  nuon api /v1/apps/{app_id} --info
  nuon api --list`,
		Args:              cobra.ArbitraryArgs,
		PersistentPreRunE: initAPI,
		RunE:              runAPI,
	}

	root.Flags().StringP("method", "X", "", "HTTP method override (GET, POST, PUT, PATCH, DELETE)")
	root.Flags().StringArrayP("query", "q", nil, "Query parameter as key=value (repeatable)")
	root.Flags().Bool("list", false, "Browse available API endpoints interactively")
	root.Flags().Bool("info", false, "Show endpoint details (params, body schema) instead of executing")
	root.Flags().Bool("raw", false, "Output raw JSON without formatting")

	root.AddCommand(tuiCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
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
		result, err := browser.Run(api, cfg.APIURL)
		if err != nil {
			return err
		}
		switch result.Action {
		case browser.ActionSelect:
			fmt.Println(result.Route.DisplayName())
			return nil
		case browser.ActionExecute:
			// Fall through to execute the GET request
			args = []string{result.Route.Path}
		default:
			return nil
		}
	}

	if len(args) == 0 {
		return cmd.Help()
	}

	path := args[0]

	// --info: show endpoint details and exit
	showInfo, _ := cmd.Flags().GetBool("info")
	if showInfo {
		routes := api.Lookup(path)
		if len(routes) == 0 {
			return fmt.Errorf("no endpoint found for path: %s", path)
		}
		output.PrintEndpointInfo(routes, cfg.APIURL)
		return nil
	}

	raw, _ := cmd.Flags().GetBool("raw")
	methodOverride, _ := cmd.Flags().GetString("method")

	var payload string
	if len(args) > 1 {
		payload = args[1]
	}

	c := client.New(cfg)

	req, err := dispatch.Resolve(api, path, payload, methodOverride, cfg, c)
	if err != nil {
		return err
	}

	// Parse -q key=value pairs into query params
	queryFlags, _ := cmd.Flags().GetStringArray("query")
	var queryParams []client.QueryParam
	for _, qf := range queryFlags {
		k, v, ok := strings.Cut(qf, "=")
		if !ok {
			return fmt.Errorf("invalid query parameter %q (expected key=value)", qf)
		}
		queryParams = append(queryParams, client.QueryParam{Key: k, Value: v})
	}

	resp, err := c.Do(req.Method, req.Path, req.Payload, queryParams...)
	if err != nil {
		return err
	}

	return output.Print(resp, raw)
}

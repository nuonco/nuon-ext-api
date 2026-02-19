package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/nuonco/nuon-ext-api/internal/config"
)

var cfg *config.Config

func Execute() {
	cfg = config.Load()

	root := &cobra.Command{
		Use:   "nuon-ext-api",
		Short: "API client for the Nuon public API",
	}

	root.AddCommand(apiCmd())
	root.AddCommand(tuiCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

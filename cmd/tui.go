package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func tuiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Interactive API explorer TUI (coming soon)",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("The TUI is not yet implemented.")
			return nil
		},
	}
}

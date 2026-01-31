package cmd

import (
	"cosoft-cli/internal/ui"
	"fmt"

	"github.com/spf13/cobra"
)

// quickBookCmd represents the quickBook command
var quickBookCmd = &cobra.Command{
	Use:     "quick-book",
	Short:   "Quickly book a meeting room",
	Long:    `Quick book allows you to quickly reserve a meeting room for a specified duration.`,
	PreRunE: requireAuth,
	Run: func(cmd *cobra.Command, args []string) {
		ui := ui.NewUI()
		if err := ui.StartApp("quick-book", false); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(quickBookCmd)
}

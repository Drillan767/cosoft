package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var landingCmd = &cobra.Command{
	Use: "landingCmd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Default command was run")
	},
}

func init() {
	rootCmd.AddCommand(landingCmd)
}

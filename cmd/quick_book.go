/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickBookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickBookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

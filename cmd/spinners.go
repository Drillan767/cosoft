/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cosoft-cli/internal/sandbox"

	"github.com/spf13/cobra"
)

// spinnersCmd represents the spinners command
var spinnersCmd = &cobra.Command{
	Use:   "spinners",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sandbox := sandbox.NewSandbox()

		sandbox.Run()
	},
}

func init() {
	rootCmd.AddCommand(spinnersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spinnersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spinnersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

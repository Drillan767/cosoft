/*
Copyright © 2026 NAME HERE <jlevarato@proton.me>
*/
package cmd

import (
	"cosoft-cli/internal/ui"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticates you to Cosoft",
	Run: func(cmd *cobra.Command, args []string) {
		ui := ui.NewUI()
		creds, err := ui.LoginFormWithLayout()
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		// Handle the credentials here
		cmd.Printf("Successfully logged in as: %s\n", creds.Email)
		// TODO: Add your authentication logic here
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

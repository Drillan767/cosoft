/*
Copyright © 2026 NAME HERE <jlevarato@proton.me>
*/
package cmd

import (
	"cosoft-cli/internal/auth"
	"cosoft-cli/internal/ui"
	"time"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticates you to Cosoft",
	Run: func(cmd *cobra.Command, args []string) {
		ui := ui.NewUI()
		loginModel, err := ui.LoginForm()
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		authService := auth.NewAuthService()

		// Calculate expiration (1 week from now)
		expirationDate := time.Now().Add(7 * 24 * time.Hour)

		if err := authService.SaveAuthData(loginModel.GetUser(), expirationDate); err != nil {
			cmd.PrintErrf("Failed to save token: %v\n", err)
			return
		}

		cmd.Printf("✓ Successfully logged in!\n")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

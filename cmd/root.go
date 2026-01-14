/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cosoft-cli/internal/auth"
	"cosoft-cli/internal/settings"
	"cosoft-cli/internal/ui"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cosoft",
	Short: "A brief description of your application",
	Long: `Cosoft CLI allows you to interact with the Cosoft's booking system without using the website.
	Through it, you can book a meeting room, list, see and cancel any reservation you've previously made.
	`,
	PreRunE: requireAuth,
	Run: func(cmd *cobra.Command, args []string) {
		ui := ui.NewUI()
		if err := ui.StartApp("landing", true); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := settings.LoadConfiguration()

		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func requireAuth(cmd *cobra.Command, args []string) error {
	authService := auth.NewAuthService()

	if authService.IsAuthenticated() {
		return nil
	}

	// Not authenticated, show form
	uiInstance := ui.NewUI()
	loginModel, err := uiInstance.LoginForm()

	if err != nil {
		return err
	}

	user := loginModel.GetUser()

	// Check if token is actually present (login succeeded)
	if user == nil || user.JwtToken == "" {
		return fmt.Errorf("authentication cancelled or failed")
	}

	// Todo: replace with actual expiration date from actual token
	expirationDate := time.Now().Add(7 * 24 * time.Hour) // 1 week

	return authService.SaveAuthData(user, expirationDate)
}

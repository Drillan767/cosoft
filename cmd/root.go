/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cosoft-cli/internal/ui"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cosoft",
	Short: "A brief description of your application",
	Long: `Cosoft CLI allows you to interact with the Cosoft's booking system without using the website.
	Through it, you can book a meeting room, list, see and cancel any reservation you've previously made.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ui := ui.NewUI()
		if err := ui.StartApp("landing", true); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check for existence of ~/.cosoft and its content
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cosoft-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "Quickly book a room with arguments",
	Run: func(cmd *cobra.Command, args []string) {
		nbUsers, err := cmd.Flags().GetInt("nbUsers")
		if err != nil {
			fmt.Println(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
		}

		date, err := cmd.Flags().GetString("time")
		if err != nil {
			fmt.Println(err)
		}

		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("reservations called")
		fmt.Println(nbUsers, name, date, duration)
	},
}

func init() {
	bookCmd.Flags().IntP("nbUsers", "u", 1, "For how many people?")
	bookCmd.Flags().StringP("name", "n", "", "If you want a room in particular")
	bookCmd.Flags().StringP("time", "t", "", "Book on a specific date")
	bookCmd.Flags().IntP("duration", "d", 30, "Duration of the booking in minutes (Must be a multiple of 15 minutes)")
	rootCmd.AddCommand(bookCmd)
}

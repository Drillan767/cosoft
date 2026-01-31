package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "Quickly book a room with arguments",
	Run: func(cmd *cobra.Command, args []string) {
		nbUsers, err := cmd.Flags().GetInt("capacity")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if nbUsers < 1 {
			nbUsers = 1
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		date, err := cmd.Flags().GetString("time")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		parsedTime, err := time.Parse("2006-01-02T15:04", date)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if parsedTime.Before(time.Now()) {
			fmt.Println("The date needs to be in the future")
			os.Exit(1)
		}

		if parsedTime.Minute()%15 != 0 {
			fmt.Println("Time needs to be rounded to a quarter")
			os.Exit(1)
		}

		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if duration == 0 || duration%15 != 0 {
			fmt.Println("Duration must be a multiple of 15")
			os.Exit(1)
		}

		fmt.Println("reservations called")
		fmt.Println(nbUsers, name, date, duration)
	},
}

func init() {
	bookCmd.Flags().IntP(
		"capacity",
		"c",
		1,
		"For how many people?",
	)

	bookCmd.Flags().StringP(
		"name",
		"n",
		"",
		"If you want a room in particular. Will pick the 1st available if not provided.",
	)

	bookCmd.Flags().StringP(
		"time",
		"d",
		"",
		"Expected format: yyyy-MM-ddTHH:mm, cannot be in the past, round the time to the closest quarter",
	)

	bookCmd.Flags().IntP(
		"duration",
		"l",
		30,
		"Duration of the booking in minutes (Must be a multiple of 15 minutes)",
	)

	rootCmd.AddCommand(bookCmd)
}

package cmd

import (
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/storage"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var roomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "List all rooms in HUB612",
	Run: func(cmd *cobra.Command, args []string) {
		cmn := common.NewCommon()

		configDir, err := os.UserConfigDir()

		if err != nil {
			log.Fatal(err)
		}

		path := fmt.Sprintf("%s/cosoft/data.db", configDir)

		store, err := storage.NewStore(path)

		if err != nil {
			log.Fatal(err)
		}

		rooms, err := store.GetRooms()

		if err != nil {
			log.Fatal(err)
		}

		headers := []string{"NAME", "CAPACITY", "PRICE"}

		rows := make([][]string, len(rooms))

		for i, room := range rooms {
			rows[i] = []string{
				room.Name,
				strconv.Itoa(room.MaxUsers) + " person(s)",
				strconv.FormatFloat(room.Price, 'g', 5, 64) + " credits",
			}
		}

		fmt.Println(cmn.CreateTable(headers, rows))
	},
}

func init() {
	rootCmd.AddCommand(roomsCmd)
}

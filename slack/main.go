package main

import (
	"cosoft-cli/internal/slackbot"
	"cosoft-cli/internal/storage"
	"log"
)

func main() {
	store, err := storage.NewStore("./slack/database.db")

	if err != nil {
		log.Fatal(err)
	}

	// Ensure database exists and is migrated
	err = store.SetupDatabase()

	if err != nil {
		log.Fatal(err)
	}

	bot := slackbot.NewBot()
	bot.StartServer()
}

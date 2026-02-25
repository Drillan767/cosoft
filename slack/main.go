package main

import (
	"cosoft-cli/internal/slackbot"
	"cosoft-cli/internal/slackbot/services"
	"cosoft-cli/internal/storage"
	"log"
	"os"
	_ "time/tzdata"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./slack/database.db"
	}

	store, err := storage.NewStore(dbPath)

	if err != nil {
		log.Fatal(err)
	}

	// Ensure database exists and is migrated
	err = store.SetupDatabase()

	if err != nil {
		log.Fatal(err)
	}

	service := services.NewSlackService(store)
	bot := slackbot.NewBot(service)
	bot.StartServer()
}

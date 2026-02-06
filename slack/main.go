package main

import "cosoft-cli/internal/slackbot"

func main() {
	bot := slackbot.NewBot()
	bot.StartServer()
}

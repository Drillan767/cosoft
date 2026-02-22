package slackbot

import "cosoft-cli/internal/slackbot/services"

type Bot struct {
	service *services.SlackService
}

func NewBot(service *services.SlackService) *Bot {
	return &Bot{service: service}
}

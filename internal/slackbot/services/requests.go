package services

import (
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
)

func (s *SlackService) ParseSlackCommand(request models.Request) (*slack.Block, error) {
	// TODO: rename this function to something more logic
	user, err := s.store.GetUserData(&request.UserId)

	if err != nil {
		return nil, err
	}

	menu := slack.MainMenu(*user)

	return &menu, nil

}

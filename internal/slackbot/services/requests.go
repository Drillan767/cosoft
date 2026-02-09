package services

import (
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"fmt"
)

func (s *SlackService) ParseSlackCommand(request models.Request) (*slack.Block, error) {
	// TODO: rename this function to something more logic
	user, err := s.store.GetUserData(&request.UserId)

	if err != nil {
		return nil, err
	}

	welcomeMessage := fmt.Sprintf(
		"Vous êtes connecté(e) en tant que *%s %s* (%s)",
		user.FirstName,
		user.LastName,
		user.Email,
	)

	creditsMessage := fmt.Sprintf("Il vous reste *%.2f* credits", user.Credits)

	return &slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewMrkDwn(welcomeMessage),
			slack.NewMrkDwn(creditsMessage),
		},
	}, nil
}

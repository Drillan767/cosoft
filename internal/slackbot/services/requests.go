package services

import (
	"bytes"
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"encoding/json"
	"net/http"
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

func (s *SlackService) UpdateMessage(responseUrl string, blocks slack.Block) error {
	blocks.ReplaceOriginal = true

	jsonPayload, err := json.Marshal(blocks)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", responseUrl, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)

	return err
}

func (s *SlackService) ShowQuickBook(action models.MenuSelection) error {
	blocks := slack.QuickBookMenu()

	return s.UpdateMessage(action.ResponseURL, blocks)
}

func (s *SlackService) ShowMainMenu(action models.MenuSelection) error {
	user, err := s.store.GetUserData(&action.User.ID)

	if err != nil {
		return err
	}

	blocks := slack.MainMenu(*user)
	return s.UpdateMessage(action.ResponseURL, blocks)
}

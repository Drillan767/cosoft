package services

import (
	"bytes"
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func (s *SlackService) DispatchModal(wrapper slack.ModalWrapper) error {
	jsonBlocks, err := json.Marshal(wrapper)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/views.open", bytes.NewBuffer(jsonBlocks))

	if err != nil {
		return err
	}

	fmt.Println(string(jsonBlocks))

	accessToken := os.Getenv("SLACK_ACCESS_TOKEN")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Println("Error dispatching modal:", buf.String())
	}

	return nil
}

func (s *SlackService) ShowQuickBook(action models.MenuSelection) error {
	modal := slack.NewQuickbook(action.ResponseURL)

	wrapper := slack.ModalWrapper{
		TriggerId: action.TriggerID,
		View:      modal,
	}

	return s.DispatchModal(wrapper)
}

func (s *SlackService) ShowMainMenu(action models.MenuSelection) error {
	user, err := s.store.GetUserData(&action.User.ID)

	if err != nil {
		return err
	}

	blocks := slack.MainMenu(*user)
	return s.UpdateMessage(action.ResponseURL, blocks)
}

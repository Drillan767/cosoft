package services

import (
	"bytes"
	"cosoft-cli/internal/slackbot/views"
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"encoding/json"
	"net/http"
)

func (s *SlackService) HandleInteraction(payload string) error {
	var result models.InteractionDiscovery

	err := json.Unmarshal([]byte(payload), &result)

	if err != nil {
		return err
	}

	dbView, err := s.store.GetSlackState(result.User.ID)

	if err != nil {
		return err
	}

	var view views.View

	if dbView != nil {
		view, err = views.RestoreView(dbView.MessageType, dbView.Payload)

		if err != nil {
			return err
		}
	}

	newView, cmd := view.Update(views.Action{
		ActionID: result.Actions[0].ActionID,
		Values:   result.State.Values,
	})

	switch c := cmd.(type) {
	case *views.LoginCmd:
		err = s.LogInUser(c.Email, c.Password, result.User.ID)

		if err != nil {
			errMsg := ":red_circle: Identifiant / mot de passe incorrect"
			if loginView, ok := newView.(*views.LoginView); ok {
				loginView.Error = &errMsg
			}
		} else {
			user, err := s.store.GetUserData(&result.User.ID)

			if err != nil {
				return err
			}

			newView = &views.LandingView{
				User: *user,
			}
		}
	}

	err = s.store.SetSlackState(result.User.ID, views.ViewType(newView), newView)

	if err != nil {
		return err
	}

	blocks := views.RenderView(newView)
	return s.SendToSlack(result.ResponseURL, blocks)
}

func (s *SlackService) SendToSlack(responseUrl string, blocks slack.Block) error {

	type BlockWrapper struct {
		ResponseType string `json:"response_type"`
		Blocks       slack.Block
	}

	jsonBlocks, err := json.Marshal(blocks)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", responseUrl, bytes.NewBuffer(jsonBlocks))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)

	return err
}

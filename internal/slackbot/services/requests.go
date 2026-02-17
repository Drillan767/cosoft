package services

import (
	"bytes"
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/slackbot/views"
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
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

		if view == nil {
			return fmt.Errorf("no active view for user %s", result.User.ID)
		}
	}

	newView, cmd := view.Update(views.Action{
		ActionID: result.Actions[0].ActionID,
		Values:   result.State.Values,
	})

	switch c := cmd.(type) {
	case nil:
		// No command and view didn't change (e.g. select events) — nothing to do
		if dbView != nil && views.ViewType(newView) == dbView.MessageType {
			return nil
		}
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
	case *views.LandingCmd:
		user, err := s.store.GetUserData(&result.User.ID)

		if err != nil {
			return err
		}

		newView = &views.LandingView{
			User: *user,
		}
	case *views.QuickBookCmd:
		rooms, err := s.getRoomAvailabilities(
			result.User.ID,
			c.NbPeople,
			c.Duration,
		)

		if err != nil {
			errMsg := ":red_circle: La réservation a échoué"
			if qbView, ok := newView.(*views.QuickBookView); ok {
				qbView.Error = &errMsg
			}
		} else {
			qbView := newView.(*views.QuickBookView)
			qbView.Phase = 2
			qbView.Rooms = &rooms

			err := s.store.SetSlackState(result.User.ID, views.ViewType(qbView), qbView)

			if err != nil {
				return err
			}

			blocks := views.RenderView(qbView)
			err = s.SendToSlack(result.ResponseURL, blocks)

			if err != nil {
				return err
			}

			room, err := s.bookRoom(
				result.User.ID,
				c.NbPeople,
				c.Duration,
				rooms,
			)

			if err != nil {
				errMsg := err.Error()
				qbView.Error = &errMsg
			} else {
				qbView.Phase = 3
				qbView.PickedRoom = room
			}

			err = s.store.SetSlackState(result.User.ID, views.ViewType(qbView), qbView)

			if err != nil {
				return err
			}

			blocks = views.RenderView(qbView)
			return s.SendToSlack(result.ResponseURL, blocks)
		}
	}

	err = s.store.SetSlackState(result.User.ID, views.ViewType(newView), newView)

	if err != nil {
		return err
	}

	blocks := views.RenderView(newView)
	return s.SendToSlack(result.ResponseURL, blocks)
}

func (s *SlackService) SetSlackState(slackUserId, messageType string, state any) error {
	return s.store.SetSlackState(slackUserId, messageType, state)
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

func (s *SlackService) getRoomAvailabilities(slackUserId string, nbPeople, duration int) ([]models.Room, error) {
	user, err := s.store.GetUserData(&slackUserId)

	if err != nil {
		return nil, err
	}

	dt := common.GetClosestQuarterHour()
	apiClient := api.NewApi()

	payload := api.CosoftAvailabilityPayload{
		DateTime: dt,
		NbPeople: nbPeople,
		Duration: duration,
	}

	rooms, err := apiClient.GetAvailableRooms(user.WAuth, user.WAuthRefresh, payload)

	if err != nil {
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, fmt.Errorf(":red_circle: Aucune salle disponible")
	}

	return rooms, nil
}

func (s *SlackService) bookRoom(slackUserId string, nbPeople, duration int, rooms []models.Room) (*models.Room, error) {
	user, err := s.store.GetUserData(&slackUserId)

	if err != nil {
		return nil, err
	}

	var pickedRoom *models.Room

	for _, room := range rooms {
		if room.NbUsers >= nbPeople {
			pickedRoom = &room
			break
		}
	}

	if pickedRoom == nil {
		return nil, fmt.Errorf(":red_circle: Aucune salle disponible")
	}

	if user.Credits < pickedRoom.Price {
		return nil, fmt.Errorf(":red_circle: Pas assez de crédits pour faire une réservation")
	}

	payload := api.CosoftBookingPayload{
		CosoftAvailabilityPayload: api.CosoftAvailabilityPayload{
			DateTime: common.GetClosestQuarterHour(),
			NbPeople: nbPeople,
			Duration: duration,
		},
		UserCredits: user.Credits,
		Room:        *pickedRoom,
	}

	apiClient := api.NewApi()

	err = apiClient.BookRoom(user.WAuth, user.WAuthRefresh, payload)

	if err != nil {
		return nil, err
	}

	return pickedRoom, nil
}

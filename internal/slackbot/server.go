package slackbot

import (
	"bytes"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type LoginFailedPayload struct {
	Password string `json:"password"`
}

type LoginFeedback struct {
	ResponseAction string              `json:"response_action"`
	Errors         *LoginFailedPayload `json:"errors,omitempty"`
}

func (b *Bot) StartServer() {
	s := http.Server{
		Addr: ":11111",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Request received", slog.String("method", r.Method), slog.String("url", r.URL.String()))

			switch r.URL.String() {
			case "/book":
				b.handleRequests(w, r)
				break
			case "/interact":
				b.handleInteractions(w, r)
			default:
				fmt.Println("Unknown URL", r.URL.String())
			}
		}),
	}

	slog.Info("Server is starting...")

	err := s.ListenAndServe()

	if err != nil {
		slog.Error("failed to listen and serve", "err", err.Error())
	}
}

func (b *Bot) handleRequests(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
		return
	}

	slackRequest := models.Request{
		Command:     r.Form.Get("command"),
		Text:        r.Form.Get("text"),
		UserId:      r.Form.Get("user_id"),
		ResponseUrl: r.Form.Get("response_url"),
		TriggerId:   r.Form.Get("trigger_id"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"response_type":"ephemeral","text":"Chargement en cours..."}`))

	go func() {

		authenticated := b.service.IsSlackAuthenticated(slackRequest)

		if !authenticated {
			b.service.DisplayLogin(slackRequest)
			return
		}

		blocks, err := b.service.ParseSlackCommand(slackRequest)

		if err != nil {
			fmt.Println(err)
			return
		}

		blocks.ResponseType = "ephemeral"

		jsonBlocks, err := json.Marshal(blocks)

		if err != nil {
			fmt.Println(err)
			return
		}

		req, err := http.NewRequest("POST", slackRequest.ResponseUrl, bytes.NewBuffer(jsonBlocks))

		if err != nil {
			fmt.Println(err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		http.DefaultClient.Do(req)
	}()
}

func (b *Bot) handleInteractions(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
		return
	}

	payload := r.Form.Get("payload")
	debug(payload)

	// Discovering what's being sent to us
	var envelope struct {
		Type string `json:"type"`
	}

	err = json.Unmarshal([]byte(payload), &envelope)

	if err != nil {
		fmt.Println(err)
		return
	}

	switch envelope.Type {
	case "view_submission":
		b.handleModalAction(payload, w)
		break
	case "block_actions":
		b.handleMenuAction(payload, w)
		break
	}
}

func (b *Bot) handleModalAction(payload string, w http.ResponseWriter) {
	type Modal struct {
		View struct {
			CallbackID string `json:"callback_id"`
		} `json:"view"`
	}

	var modal Modal

	err := json.Unmarshal([]byte(payload), &modal)

	if err != nil {
		fmt.Println(err)
		return
	}

	switch modal.View.CallbackID {
	case "login_modal":
		b.handleLoginModal(payload, w)
	case "quickbook_modal":
		b.handleQuickbookModal(payload, w)
	}
}

func (b *Bot) handleLoginModal(payload string, w http.ResponseWriter) {
	var viewResponse models.SlackLoginResponse

	err := json.Unmarshal([]byte(payload), &viewResponse)

	if err != nil {
		fmt.Println(err)
		return
	}

	// PTSD from Drupal 8's forms overriding.
	email := viewResponse.View.State.Values.Email.Email.Value
	password := viewResponse.View.State.Values.Password.Password.Value
	responseUrl := viewResponse.View.PrivateMetadata
	slackUserId := viewResponse.User.ID

	err = b.service.LogInUser(email, password, slackUserId, responseUrl)

	if err != nil {
		feedback := LoginFeedback{
			ResponseAction: "errors",
			Errors: &LoginFailedPayload{
				Password: err.Error(),
			},
		}

		jsonValue, err := json.Marshal(feedback)

		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonValue)

		fmt.Println(err)
		return
	}

	feedback := LoginFeedback{
		ResponseAction: "clear",
	}

	jsonValue, err := json.Marshal(feedback)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonValue)
	return
}

func (b *Bot) handleQuickbookModal(payload string, w http.ResponseWriter) {
	var viewResponse models.SlackQuickBookResponse

	err := json.Unmarshal([]byte(payload), &viewResponse)

	if err != nil {
		fmt.Println(err)
		return
	}

	// lmao.
	duration := viewResponse.View.State.Values.Duration.Duration.SelectedOption.Value
	nbPeople := viewResponse.View.State.Values.NbPeople.NbPeople.SelectedOption.Value

	fmt.Println(duration, nbPeople)
}

func (b *Bot) handleMenuAction(payload string, w http.ResponseWriter) {
	var action models.MenuSelection

	err := json.Unmarshal([]byte(payload), &action)

	if err != nil {
		fmt.Println(err)
		return
	}

	actionName := action.Actions[0].ActionID

	switch actionName {
	case "main-menu":
		w.WriteHeader(http.StatusOK)
		go func() {
			if err := b.service.ShowMainMenu(action); err != nil {
				fmt.Println(err)
			}
		}()
	case "quick-book":
		if err := b.service.ShowQuickBook(action); err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func debug(payload interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	file.WriteString(fmt.Sprintf("%v \n\n", payload))
}

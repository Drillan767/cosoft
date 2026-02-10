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

	w.WriteHeader(http.StatusOK)

	go func() {
		s := b.service

		authenticated := s.IsSlackAuthenticated(slackRequest)

		if !authenticated {
			s.DisplayLogin(slackRequest)
			return
		}

		blocks, err := s.ParseSlackCommand(slackRequest)

		if err != nil {
			fmt.Println(err)
			return
		}

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

	var viewResponse models.SlackLoginResponse

	err = json.Unmarshal([]byte(payload), &viewResponse)

	if err != nil {
		fmt.Println(err)
		return
	}

	s, err := services.NewSlackService()

	if err != nil {
		fmt.Println(err)
		return
	}

	s := b.service

	// PTSD from Drupal 8's forms overriding.
	email := viewResponse.View.State.Values.Email.Email.Value
	password := viewResponse.View.State.Values.Password.Password.Value
	responseUrl := viewResponse.View.PrivateMetadata
	slackUserId := viewResponse.User.ID

	err = s.LogInUser(email, password, slackUserId, responseUrl)

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

func debug(payload interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	file.WriteString(fmt.Sprintf("%v \n\n", payload))
}

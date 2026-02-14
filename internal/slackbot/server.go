package slackbot

import (
	"cosoft-cli/internal/slackbot/views"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

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

	// Clear out user's old slack states
	err = b.service.ClearUserStates(slackRequest)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"response_type":"ephemeral","text":"Chargement en cours..."}`))

	go func() {
		view, err := b.service.AuthGuard(slackRequest)

		if err != nil {
			fmt.Println(err)
			return
		}

		if view != nil {
			blocks := views.RenderView(view)
			err = b.service.SendToSlack(slackRequest.ResponseUrl, blocks)

			if err != nil {
				fmt.Println(err)
				return
			}
		}

		mainMenu := views.LandingView{}

		blocks := views.RenderView(&mainMenu)

		err = b.service.SendToSlack(slackRequest.ResponseUrl, blocks)

		if err != nil {
			fmt.Println(err)
			return
		}
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
	case "block_actions":
		// TODO:
		// - Read view from DB using slack_message_id
		// - view = view.Update
		// - Write view to DB
		// - RenderView() => slack.Block
		// - Send to Slack
		// b.handleMenuAction(payload, w)
		break
	}
}

func debug(payload interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	file.WriteString(fmt.Sprintf("%v \n\n", payload))
}

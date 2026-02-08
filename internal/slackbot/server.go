package slackbot

import (
	"cosoft-cli/internal/slackbot/services"
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

			err := r.ParseForm()

			if err != nil {
				fmt.Println(err)
			}

			slackRequest := models.Request{
				Command:     r.Form.Get("command"),
				Text:        r.Form.Get("text"),
				UserId:      r.Form.Get("user_id"),
				ResponseUrl: r.Form.Get("response_url"),
				TriggerId:   r.Form.Get("trigger_id"),
			}

			s, err := services.NewSlackService()

			if err != nil {
				fmt.Println(err)
				w.Write([]byte(err.Error()))
				return
			}

			authenticated := s.IsSlackAuthenticated(slackRequest)

			if !authenticated {
				s.DisplayLogin(slackRequest)
				return
			}

			blocks := services.ParseSlackCommand(slackRequest)

			jsonBlocks, err := json.Marshal(blocks)

			if err != nil {
				fmt.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonBlocks)

		}),
	}

	slog.Info("Server is starting...")

	err := s.ListenAndServe()

	if err != nil {
		slog.Error("failed to listen and serve", "err", err.Error())
	}
}

func debug(payload interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	file.WriteString(fmt.Sprintf("%v \n\n", payload))
}

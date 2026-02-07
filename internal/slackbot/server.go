package slackbot

import (
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

			slackRequest := Request{
				Command:     r.Form.Get("command"),
				Text:        r.Form.Get("text"),
				UserId:      r.Form.Get("user_id"),
				ResponseUrl: r.Form.Get("response_url"),
			}

			debug(slackRequest)

			fmt.Println(slackRequest)
			w.Write([]byte(slackRequest.ResponseUrl))

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

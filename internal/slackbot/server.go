package slackbot

import (
	"cosoft-cli/internal/slackbot/services"
	"cosoft-cli/shared/models"
	"encoding/json"
	"fmt"
	"io"
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
}

/** EXAMPLE RESPONSE PAYLOAD
{
	"type":"view_submission",
	"team":{
		"id":"T0ADB9JET2L",
		"domain":"cosoftlamaudite"
	},
	"user":{
		"id":"U0ACAJPPR7D",
		"username":"joseph",
		"name":"joseph",
		"team_id":"T0ADB9JET2L"
	},
	"api_app_id":"A0ACEV21BN2",
	"token":"J3r5ymWUxW5vyU7w66kQYtEF",
	"trigger_id":"10479377324721.10453324503088.b25873bd2eb05eb5303316f8fa763a46",
	"view":{
		"id":"V0ADNBLJFRU",
		"team_id":"T0ADB9JET2L",
		"type":"modal",
		"blocks":[
			{
				"type":"input",
				"block_id":"email",
				"label":{
					"type":"plain_text",
					"text":"Email",
					"emoji":true
				},
				"optional":false,
				"dispatch_action":false,
				"element":{
					"type":"plain_text_input",
					"action_id":"email",
					"dispatch_action_config":{
						"trigger_actions_on":[
							"on_enter_pressed"
						]
					}
				}
			},
			{
				"type":"input",
				"block_id":"password",
				"label":{
					"type":"plain_text",
					"text":"Mot de passe",
					"emoji":true
				},
				"optional":false,
				"dispatch_action":false,
				"element":{
					"type":"plain_text_input",
					"action_id":"password",
					"dispatch_action_config":{
						"trigger_actions_on":[
							"on_enter_pressed"
						]
					}
				}
			},
			{
				"type":"context",
				"block_id":"uFpGL",
				"elements":[
					{
						"type":"plain_text",
						"text":":warning: Le mot de passe est affich\u00e9 en clair dans le champ",
						"emoji":true
					}
				]
			}
		],
		"private_metadata":"",
		"callback_id":"login_modal",
		"state":{
			"values":{
				"email":{
					"email":{
						"type":"plain_text_input",
						"value":"asdsadasd"
					}
				},
				"password":{
					"password":{
						"type":"plain_text_input",
						"value":"rgfdfgdfgdfgd"
					}
				}
			}
		},
		"hash":"1770556275.Ls40d2cz",
		"title":{
			"type":"plain_text",
			"text":"Connexion",
			"emoji":true
		},
		"clear_on_close":false,
		"notify_on_close":false,
		"close":{
			"type":"plain_text",
			"text":"Fermer",
			"emoji":true
		},
		"submit":{
			"type":"plain_text",
			"text":"Connexion",
			"emoji":true
		},
		"previous_view_id":null,
		"root_view_id":"V0ADNBLJFRU",
		"app_id":"A0ACEV21BN2",
		"external_id":"",
		"app_installed_team_id":"T0ADB9JET2L",
		"bot_id":"B0ACDHD0Y21"
	},
	"response_urls":[

	],
	"is_enterprise_install":false,
	"enterprise":null
}
*/

func (b *Bot) handleInteractions(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	debug(string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("OK"))
	return
}

func debug(payload interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	file.WriteString(fmt.Sprintf("%v \n\n", payload))
}

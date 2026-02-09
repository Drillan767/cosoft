package models

import (
	"github.com/charmbracelet/bubbles/spinner"
)

type GlobalState struct {
	currentPage string
	spinner     spinner.Model
	loading     bool
	quickAction bool
}

type Selection struct {
	Choice string
}

type Room struct {
	Id      string
	Name    string
	NbUsers int
	Price   float64
}

type UnavailableSlot struct {
	Title string `json:"Title"`
	Start string `json:"Start"`
	End   string `json:"End"`
}

type RoomUsage struct {
	Name      string
	Id        string
	UsedSlots []UnavailableSlot
}

type Request struct {
	UserId      string
	Command     string
	Text        string
	ResponseUrl string
	TriggerId   string
}

type SlackLoginResponse struct {
	Type string `json:"type"`
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		TeamID   string `json:"team_id"`
	} `json:"user"`
	APIAppID  string `json:"api_app_id"`
	Token     string `json:"token"`
	TriggerID string `json:"trigger_id"`
	View      struct {
		ID              string `json:"id"`
		TeamID          string `json:"team_id"`
		Type            string `json:"type"`
		PrivateMetadata string `json:"private_metadata"`
		CallbackID      string `json:"callback_id"`
		State           struct {
			Values struct {
				Email struct {
					Email struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"email"`
				} `json:"email"`
				Password struct {
					Password struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"password"`
				} `json:"password"`
			} `json:"values"`
		} `json:"state"`
		Hash               string `json:"hash"`
		PreviousViewID     any    `json:"previous_view_id"`
		RootViewID         string `json:"root_view_id"`
		AppID              string `json:"app_id"`
		ExternalID         string `json:"external_id"`
		AppInstalledTeamID string `json:"app_installed_team_id"`
		BotID              string `json:"bot_id"`
	} `json:"view"`
	ResponseUrls        []any `json:"response_urls"`
	IsEnterpriseInstall bool  `json:"is_enterprise_install"`
	Enterprise          any   `json:"enterprise"`
}

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

type SlackQuickBookResponse struct {
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
				Duration struct {
					Duration struct {
						Type           string `json:"type"`
						SelectedOption struct {
							Text struct {
								Type  string `json:"type"`
								Text  string `json:"text"`
								Emoji bool   `json:"emoji"`
							} `json:"text"`
							Value string `json:"value"`
						} `json:"selected_option"`
					} `json:"duration"`
				} `json:"duration"`
				NbPeople struct {
					NbPeople struct {
						Type           string `json:"type"`
						SelectedOption struct {
							Text struct {
								Type  string `json:"type"`
								Text  string `json:"text"`
								Emoji bool   `json:"emoji"`
							} `json:"text"`
							Value string `json:"value"`
						} `json:"selected_option"`
					} `json:"nbPeople"`
				} `json:"nbPeople"`
			} `json:"values"`
		} `json:"state"`
		Hash string `json:"hash"`
	} `json:"view"`
}

type MenuSelection struct {
	Type string `json:"type"`
	User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		TeamID   string `json:"team_id"`
	} `json:"user"`
	APIAppID  string `json:"api_app_id"`
	Token     string `json:"token"`
	Container struct {
		Type        string `json:"type"`
		MessageTs   string `json:"message_ts"`
		ChannelID   string `json:"channel_id"`
		IsEphemeral bool   `json:"is_ephemeral"`
	} `json:"container"`
	TriggerID string `json:"trigger_id"`
	Channel   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	ResponseURL string `json:"response_url"`
	Actions     []struct {
		ActionID string `json:"action_id"`
		BlockID  string `json:"block_id"`
		Text     struct {
			Type  string `json:"type"`
			Text  string `json:"text"`
			Emoji bool   `json:"emoji"`
		} `json:"text"`
		Value    string `json:"value"`
		Type     string `json:"type"`
		ActionTs string `json:"action_ts"`
	} `json:"actions"`
}

package models

import (
	"encoding/json"

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
	Image   string
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

type InteractionDiscovery struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	State struct {
		Values json.RawMessage `json:"values"`
	} `json:"state"`
	ResponseURL string `json:"response_url"`
	Actions     []struct {
		ActionID string `json:"action_id"`
	}
}

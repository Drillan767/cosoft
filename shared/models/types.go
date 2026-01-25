package models

import (
	"time"

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
	StartDate time.Time
	EndDate   time.Time
}

type RoomUsage struct {
	Name      string
	Id        string
	UsedSlots []UnavailableSlot
}

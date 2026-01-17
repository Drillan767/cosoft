package models

import "github.com/charmbracelet/bubbles/spinner"

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

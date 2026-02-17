package views

import (
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
	"strings"
)

type BrowseView struct {
	Phase      int
	NbPeople   int
	Duration   int
	Date       string
	Time       string
	Rooms      *[]models.Room
	PickedRoom *models.Room
	Error      *string
}

type BrowseCmd struct{}

func (b *BrowseView) Update(action Action) (View, Cmd) {

	if action.ActionID == "cancel" {
		return b, &LandingCmd{}
	} else if action.ActionID == "browse" {
		// Get room availabilities
	} else if strings.HasPrefix(action.ActionID, "book-") {
		// A room has been picked
		// Return this for now.
		return b, nil
	}

	return b, nil
}

func RenderBrowseView(b *BrowseView) slack.Block {
	return slack.BrowseMenu()
}

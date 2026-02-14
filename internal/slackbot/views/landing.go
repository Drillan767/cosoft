package views

import (
	"cosoft-cli/internal/ui/slack"
)

type LandingView struct {
	choice string
}

func (lv *LandingView) Update(action Action) (View, Cmd) {

	switch action.ActionID {
	case "quick-book":
		return lv, nil

	default:
		return lv, nil
	}
}

func RenderLandingView(lv *LandingView) slack.Block {
	return slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewHeader("Menu principal"),
			slack.NewMenuItem(
				"*Réservation rapide*\nRéserver immédiatement une salle de réunion",
				"Accéder",
				"quick-book",
			),
		},
	}
}

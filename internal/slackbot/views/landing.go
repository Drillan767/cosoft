package views

import (
	"cosoft-cli/internal/ui/slack"
	"encoding/json"
)

type LandingView struct {
	choice string
}

func (lv *LandingView) Update(state json.RawMessage) View {
	lv.choice = "quick_book"

	return lv

	/*
		Suivant l'action qui a été cliquée, retourner la vue correspondante
		- Si quickbook, renvoyer un future QuickBookView{}
	*/
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

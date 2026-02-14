package views

import "cosoft-cli/internal/ui/slack"

type QuickBookView struct{}

func (qb *QuickBookView) Update(action Action) (View, Cmd) {
	// Pending handling selection

	switch action.ActionID {
	case "cancel":
		return qb, &LandingCmd{}
	}
	return qb, nil
}

func RenderQuickBookView() slack.Block {
	return slack.QuickBookMenu()
}

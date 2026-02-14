package views

import (
	"cosoft-cli/internal/storage"
	"cosoft-cli/internal/ui/slack"
)

type LandingView struct {
	User storage.User
}

type LandingCmd struct{}

func (lv *LandingView) Update(action Action) (View, Cmd) {

	switch action.ActionID {
	case "quick-book":
		return &QuickBookView{}, nil

	default:
		return lv, nil
	}
}

func RenderLandingView(lv *LandingView) slack.Block {
	return slack.MainMenu(lv.User)
}

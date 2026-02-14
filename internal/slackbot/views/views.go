package views

import (
	"cosoft-cli/internal/ui/slack"
	"encoding/json"
)

type slackInput struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type View interface {
	Update(state json.RawMessage) View
}

func RenderView(v View) slack.Block {
	switch v := v.(type) {
	case *LoginView:
		return RenderLoginView(v)
	case *LandingView:
		return RenderLandingView(v)
	default:
		return slack.Block{}
	}
}

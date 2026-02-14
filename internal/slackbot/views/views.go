package views

import (
	"cosoft-cli/internal/ui/slack"
	"encoding/json"
	"fmt"
)

type Action struct {
	ActionID string          `json:"action_id"`
	Values   json.RawMessage `json:"values"`
}

type View interface {
	Update(action Action) (View, Cmd)
}

type Cmd interface{}

func RestoreView(messageType string, payload []byte) (View, error) {
	var view View

	switch messageType {
	case "login":
		view = &LoginView{}
		break
	case "landing":
		view = &LandingView{}
		break
	default:
		return nil, fmt.Errorf("unknown view type: %s", messageType)
	}

	err := json.Unmarshal(payload, view)

	if err != nil {
		return nil, err
	}

	return view, nil
}

func ViewType(v View) string {
	switch v.(type) {
	case *LoginView:
		return "login"
	case *LandingView:
		return "landing"
	default:
		return "unknown"
	}
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

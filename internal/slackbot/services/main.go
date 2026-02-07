package services

import (
	"cosoft-cli/internal/ui/slack"
	"cosoft-cli/shared/models"
)

func ParseSlackCommand(request models.Request) slack.Block {
	switch request.Text {
	case "book":
		return slack.Block{
			Blocks: []slack.BlockElement{
				slack.NewMrkDwn("Bro wants a *meeting room*"),
			},
		}
	default:
		return slack.Block{
			Blocks: []slack.BlockElement{
				slack.NewMrkDwn("Bro doesn't want *anything*"),
			},
		}
	}
}

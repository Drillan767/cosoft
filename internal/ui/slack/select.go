package slack

import "fmt"

type Select struct {
	Type      string          `json:"type"`
	Text      OptionContent   `json:"text"`
	Accessory OptionAccessory `json:"accessory"`
}
type OptionContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SelectOptions struct {
	OptionContent BlockPayload `json:"text"`
	Value         string       `json:"value"`
}
type OptionAccessory struct {
	Type        string          `json:"type"`
	Placeholder BlockPayload    `json:"placeholder"`
	Options     []SelectOptions `json:"options"`
	ActionID    string          `json:"action_id"`
}

func (Select) blockElement() {}

func NewSelect(
	label, placeholder, name string,
	choices []ChoicePayload,
) Select {
	options := make([]SelectOptions, len(choices))

	for i, choice := range choices {
		options[i] = SelectOptions{
			Value: choice.Value,
			OptionContent: BlockPayload{
				Type:  "plain_text",
				Text:  choice.Text,
				Emoji: true,
			},
		}
	}

	return Select{
		Type: "section",
		Text: OptionContent{
			Type: "plain_text",
			Text: fmt.Sprintf("*%s*", label),
		},
		Accessory: OptionAccessory{
			Type: "static_select",
			Placeholder: BlockPayload{
				Type:  "plain_text",
				Text:  placeholder,
				Emoji: true,
			},
			Options: options,
		},
	}

}

package slack

type RadioLabel struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type RadioItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type RadioOption struct {
	Text  RadioItem `json:"text"`
	Value string    `json:"value"`
}

type RadioElement struct {
	Type     string        `json:"type"`
	ActionID string        `json:"action_id"`
	Options  []RadioOption `json:"options"`
}

type Radio struct {
	Type    string       `json:"type"`
	BlockId string       `json:"block_id"`
	Label   RadioLabel   `json:"label"`
	Element RadioElement `json:"element"`
}

func (Radio) blockElement() {}

func NewRadio(label, name string, choices []ChoicePayload) Radio {
	options := make([]RadioOption, len(choices))

	for i, choice := range choices {
		options[i] = RadioOption{
			Text: RadioItem{
				Text: choice.Text,
				Type: "plain_text",
			},
			Value: choice.Value,
		}
	}

	return Radio{
		Type:    "input",
		BlockId: name,
		Label:   RadioLabel{"plain_text", label},
		Element: RadioElement{
			Type:     "radio_buttons",
			ActionID: name,
			Options:  options,
		},
	}
}

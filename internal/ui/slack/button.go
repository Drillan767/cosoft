package slack

type ButtonPayload struct {
	Type     string       `json:"type"`
	Text     BlockPayload `json:"text"`
	Value    string       `json:"value"`
	ActionId string       `json:"action_id"`
}

type Button struct {
	Type     string          `json:"type"`
	Elements []ButtonPayload `json:"elements"`
}

type ButtonBlockPayload struct {
	Text  string
	Value string
}

func (Button) blockElement() {}

func NewButtons(buttons []ButtonBlockPayload) Button {
	elements := make([]ButtonPayload, len(buttons))

	for i, button := range buttons {
		elements[i] = ButtonPayload{
			Type: "button",
			Text: BlockPayload{
				Text:  button.Text,
				Type:  "plain_text",
				Emoji: true,
			},
			Value:    button.Value,
			ActionId: button.Value,
		}
	}

	return Button{
		Type:     "actions",
		Elements: elements,
	}
}

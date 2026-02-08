package slack

type InputPayload struct {
	Type     string `json:"type"`
	ActionId string `json:"action_id"`
}

type LabelPayload struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Input struct {
	Type     string       `json:"type"`
	Element  InputPayload `json:"element"`
	Label    LabelPayload `json:"label"`
	Optional bool         `json:"optional"`
}

func (i Input) blockElement() {}

func NewInput(label, name string) Input {
	return Input{
		Type: "input",
		Element: InputPayload{
			Type:     "plain_text_input",
			ActionId: name,
		},
		Label: LabelPayload{
			Type: "plain_text",
			Text: label,
		},
		Optional: false,
	}
}

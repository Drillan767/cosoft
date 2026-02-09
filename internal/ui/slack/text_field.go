package slack

type InputPayload struct {
	Type     string `json:"type"`
	ActionId string `json:"action_id"`
}

type Input struct {
	Type     string       `json:"type"`
	BlockId  string       `json:"block_id"`
	Element  InputPayload `json:"element"`
	Label    BlockPayload `json:"label"`
	Optional bool         `json:"optional"`
}

func (i Input) blockElement() {}

func NewInput(label, name string) Input {
	return Input{
		Type:    "input",
		BlockId: name,
		Element: InputPayload{
			Type:     "plain_text_input",
			ActionId: name,
		},
		Label: BlockPayload{
			Type:  "plain_text",
			Text:  label,
			Emoji: true,
		},
		Optional: false,
	}
}

package slack

type MrkDwnPayload struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Mrkdwn struct {
	Type string        `json:"type"`
	Text MrkDwnPayload `json:"text"`
}

func (Mrkdwn) blockElement() {}

func NewMrkDwn(text string) Mrkdwn {
	return Mrkdwn{
		Type: "section",
		Text: MrkDwnPayload{
			Type: "mrkdwn",
			Text: text,
		},
	}
}

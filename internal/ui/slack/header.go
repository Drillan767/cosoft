package slack

type HeaderPayload struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Header struct {
	Type string        `json:"type"`
	Text HeaderPayload `json:"text"`
}

func (Header) blockElement() {}

func NewHeader(text string) Header {
	return Header{
		Type: "header",
		Text: HeaderPayload{
			Type: "plain_text",
			Text: text,
		},
	}
}

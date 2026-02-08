package slack

type Context struct {
	Type     string         `json:"type"`
	Elements []BlockPayload `json:"elements"`
}

func (c Context) blockElement() {}

func NewContext(text string) Context {
	return Context{
		Type: "context",
		Elements: []BlockPayload{
			{
				Type:  "plain_text",
				Text:  text,
				Emoji: true,
			},
		},
	}
}

package slack

type MultiMarkdown struct {
	Type   string         `json:"type"`
	Fields []BlockPayload `json:"fields"`
}

func (MultiMarkdown) blockElement() {}

func NewMultiMarkdown(texts []string) MultiMarkdown {
	fields := make([]BlockPayload, len(texts))

	for i, text := range texts {
		fields[i] = BlockPayload{
			Type: "mrkdwn",
			Text: text,
		}
	}

	return MultiMarkdown{
		Type:   "section",
		Fields: fields,
	}
}

package slack

type KitchenSink struct {
	Type     string               `json:"type"`
	Elements []KitchenSinkPayload `json:"elements"`
}
type KitchenSinkPayload struct {
	Type     string         `json:"type"`
	Elements []BlockPayload `json:"elements"`
}

func (KitchenSink) blockElement() {}

func NewKitchenSink(text string) KitchenSink {
	return KitchenSink{
		Type: "rich_text",
		Elements: []KitchenSinkPayload{{
			Type: "rich_text_preformatted",
			Elements: []BlockPayload{{
				Type: "text",
				Text: text,
			}},
		}},
	}
}

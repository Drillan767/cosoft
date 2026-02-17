package slack

type _ButtonPayload struct {
	Type     string       `json:"type"`
	Text     BlockPayload `json:"text"`
	Value    string       `json:"value"`
	ActionId string       `json:"action_id"`
	Style    *string      `json:"style,omitempty"`
}

type ActionPayload struct {
	Type     string          `json:"type"`
	Elements []ActionPayload `json:"elements"`
}

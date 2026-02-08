package slack

type ButtonContentPayload struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Emoji bool   `json:"emoji"`
}

type ButtonPayload struct {
	Type     string               `json:"type"`
	Text     ButtonContentPayload `json:"text"`
	Value    string               `json:"value"`
	ActionId string               `json:"action_id"`
	Style    *string              `json:"style,omitempty"`
}

type ActionPayload struct {
	Type     string          `json:"type"`
	Elements []ActionPayload `json:"elements"`
}

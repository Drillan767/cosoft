package slack

type ButtonPayload struct {
	Type     string       `json:"type"`
	Text     BlockPayload `json:"text"`
	Value    string       `json:"value"`
	ActionId string       `json:"action_id"`
}

type MenuItem struct {
	Mrkdwn
	Accessory ButtonPayload `json:"accessory"`
}

func (MenuItem) blockElement() {}

func NewMenuItem(text, btnText, value, actionId string) MenuItem {
	return MenuItem{
		Mrkdwn{
			Type: "section",
			Text: MrkDwnPayload{
				Text: text,
				Type: "mrkdwn",
			},
		},
		ButtonPayload{
			Type: "button",
			Text: BlockPayload{
				Type:  "plain_text",
				Text:  btnText,
				Emoji: true,
			},
			ActionId: actionId,
			Value:    value,
		},
	}
}

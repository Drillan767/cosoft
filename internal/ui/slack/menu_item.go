package slack

type MenuItem struct {
	Mrkdwn
	Accessory ButtonPayload `json:"accessory"`
}

func (MenuItem) blockElement() {}

func NewMenuItem(text, btnText, value string) MenuItem {
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
			ActionId: value,
			Value:    value,
		},
	}
}

package slack

import (
	"cosoft-cli/internal/common"
)

type TimePickerAccessory struct {
	Type        string       `json:"type"`
	InitialTime string       `json:"initial_time"`
	ActionID    string       `json:"action_id"`
	Placeholder BlockPayload `json:"placeholder"`
}

type TimePicker struct {
	Type      string              `json:"type"`
	Text      BlockPayload        `json:"text"`
	BlockId   string              `json:"block_id"`
	Accessory TimePickerAccessory `json:"accessory"`
}

func (TimePicker) blockElement() {}

func NewTimePicker(label, name, placeholder string) TimePicker {
	time := common.GetClosestQuarterHour()

	return TimePicker{
		Type:    "section",
		BlockId: name,
		Text: BlockPayload{
			Type: "mrkdwn",
			Text: label,
		},
		Accessory: TimePickerAccessory{
			Type:        "timepicker",
			InitialTime: time.Format("15:04"),
			Placeholder: BlockPayload{
				Type:  "plain_text",
				Text:  placeholder,
				Emoji: true,
			},
			ActionID: name,
		},
	}
}

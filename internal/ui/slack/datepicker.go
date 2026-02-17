package slack

import (
	"fmt"
	"time"
)

type DatePickerAccessory struct {
	Type        string       `json:"type"`
	InitialDate string       `json:"initial_date"`
	ActionID    string       `json:"action_id"`
	Placeholder BlockPayload `json:"placeholder"`
}

type DatePicker struct {
	Type      string              `json:"type"`
	Text      BlockPayload        `json:"text"`
	Accessory DatePickerAccessory `json:"accessory"`
}

func (DatePicker) blockElement() {}

func NewDatePicker(label, name, placeholder string) DatePicker {
	today := time.Now().Format(time.DateOnly)

	return DatePicker{
		Type: "section",
		Text: BlockPayload{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*%s*", label),
		},
		Accessory: DatePickerAccessory{
			Type:        "datepicker",
			InitialDate: today,
			Placeholder: BlockPayload{
				Type: "plain_text",
				Text: placeholder,
			},
			ActionID: name,
		},
	}
}

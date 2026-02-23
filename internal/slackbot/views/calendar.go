package views

import (
	"cosoft-cli/internal/slackbot/models"
	"cosoft-cli/internal/ui/slack"
	"fmt"
	"time"
)

func RenderCalendarView(s *models.CalendarState) slack.Block {
	dt := s.CurrentDate.Format("02/01/2006")
	isToday := s.CurrentDate.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour))
	actions := []slack.ChoicePayload{{Text: "Jour suivant", Value: "next-day"}}

	if !isToday {
		actions = append(
			[]slack.ChoicePayload{{Text: "Jour précédent", Value: "prev-day"}},
			actions...,
		)
	}

	return slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewHeader("Calendrier"),
			slack.NewMrkDwn(fmt.Sprintf("*%s*", dt)),
			slack.NewDivider(),
			slack.NewKitchenSink(s.Calendar),
			slack.NewMrkDwn("`█`: Créneau réservé par vous"),
			slack.NewMrkDwn("`░`: Créneau réservé par quelqu'un d'autre"),
			slack.NewButtons(actions),
			slack.NewDivider(),
			slack.NewButtons([]slack.ChoicePayload{{Text: "Retour", Value: "cancel"}}),
		},
	}
}

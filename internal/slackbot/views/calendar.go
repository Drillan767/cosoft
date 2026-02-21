package views

import (
	"cosoft-cli/internal/ui/slack"
	"fmt"
	"time"
)

type CalendarView struct {
	CurrentDate time.Time
	Calendar    string
	Error       *string
}

type CalendarCmd struct {
	Time time.Time
}

func NewCalendarView() *CalendarView {
	return &CalendarView{
		CurrentDate: time.Now(),
	}
}

func (c *CalendarView) Update(action Action) (View, Cmd) {
	switch action.ActionID {
	case "cancel":
		return c, &LandingCmd{}
	case "next-day":
		c.CurrentDate = c.CurrentDate.Add(24 * time.Hour)
		return c, &CalendarCmd{
			Time: c.CurrentDate,
		}
	case "prev-day":
		today := time.Now().Truncate(24 * time.Hour)
		if c.CurrentDate.Truncate(24 * time.Hour).Equal(today) {
			// already on today, can't go back
			return c, nil
		}
		c.CurrentDate = c.CurrentDate.Add(-24 * time.Hour)
		return c, &CalendarCmd{}
	default:
		return c, nil
	}
}

func RenderCalendarView(c *CalendarView) slack.Block {
	dt := c.CurrentDate.Format("02/01/2006")
	isToday := c.CurrentDate.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour))
	actions := []slack.ChoicePayload{{"Jour suivant", "next-day"}}

	if !isToday {
		actions = append(
			[]slack.ChoicePayload{{"Jour précédent", "prev-day"}},
			actions...,
		)
	}

	return slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewHeader("Calendrier"),
			slack.NewMrkDwn(fmt.Sprintf("*%s*", dt)),
			slack.NewDivider(),
			slack.NewKitchenSink(c.Calendar),
			slack.NewButtons(actions),
			slack.NewDivider(),
			slack.NewButtons([]slack.ChoicePayload{{"Retour", "cancel"}}),
		},
	}
}

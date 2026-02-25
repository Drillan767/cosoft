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

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
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
		if sameDay(c.CurrentDate, time.Now()) {
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
	isToday := sameDay(c.CurrentDate, time.Now())
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
			slack.NewKitchenSink(c.Calendar),
			slack.NewMrkDwn("`█`: Créneau réservé par vous"),
			slack.NewMrkDwn("`░`: Créneau réservé par quelqu'un d'autre"),
			slack.NewButtons(actions),
			slack.NewDivider(),
			slack.NewButtons([]slack.ChoicePayload{{Text: "Retour", Value: "cancel"}}),
		},
	}
}

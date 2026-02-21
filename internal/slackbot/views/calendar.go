package views

import (
	"cosoft-cli/internal/ui/slack"
	"fmt"
	"time"
)

type CalendarView struct {
	CurrentDate time.Time
	Calendar    string
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
	if action.ActionID == "cancel" {
		return c, &LandingView{}
	}

	return c, nil
}

func RenderCalendarView(c *CalendarView) slack.Block {
	dt := c.CurrentDate.Format("02/01/2006")
	isToday := c.CurrentDate.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour))
	actions := []slack.ChoicePayload{{"Jour suivant", "next"}}

	if !isToday {
		actions = append(
			[]slack.ChoicePayload{{"Jour précédent", "prev"}},
			actions...,
		)
	}

	return slack.Block{
		Blocks: []slack.BlockElement{
			slack.NewHeader("Calendrier"),
			slack.NewMrkDwn(fmt.Sprintf("*%s*", dt)),
			slack.NewDivider(),
			slack.NewMrkDwn("Ici y'aura le calendrier"),
			slack.NewButtons(actions),
			slack.NewDivider(),
			slack.NewButtons([]slack.ChoicePayload{{"Retour", "cancel"}}),
		},
	}
}

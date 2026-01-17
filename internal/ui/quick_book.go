package ui

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/ui/components"
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type QuickBookModel struct {
	// 0 select duration, 1 select nb people, 2 display future loader, 3 display results
	phase          int
	pickedDuration int
	pickedNbPeople int
	durationList   components.ListModel
	nbPeopleList   components.ListModel
	payload        *api.QuickBookPayload
}

func NewQuickBookModel() *QuickBookModel {
	selection := &api.QuickBookPayload{}

	t := getClosestQuarterHour()

	durations := []components.Item{
		{
			Label:    "30 minutes",
			Subtitle: fmt.Sprintf("From %s to %s", t.Format("15:04"), t.Add(30*time.Minute).Format("15:04")),
			Value:    30,
		},
		{
			Label:    "1 hour",
			Subtitle: fmt.Sprintf("From %s to %s", t.Format("15:04"), t.Add(60*time.Minute).Format("15:04")),
			Value:    60,
		},
	}

	peoples := []components.Item{
		{
			Label:    "1 person",
			Subtitle: "The research will include callboxes",
			Value:    1,
		},
		{
			Label:    "2 persons or more",
			Subtitle: "The research will default to classic meeting rooms",
			Value:    2,
		},
	}

	return &QuickBookModel{
		phase:          0,
		pickedDuration: 0,
		pickedNbPeople: 0,
		durationList:   components.NewListModel(durations, "For how long?"),
		nbPeopleList:   components.NewListModel(peoples, "For how many people?"),
		payload:        selection,
	}
}

func (qb *QuickBookModel) Init() tea.Cmd {
	return nil
}

func (qb *QuickBookModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Let the list component process the message first
	switch qb.phase {
	case 0:
		qb.durationList, cmd = qb.durationList.Update(msg)
	case 1:
		qb.nbPeopleList, cmd = qb.nbPeopleList.Update(msg)
	}

	// Then check for phase transitions after selection is confirmed
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			switch qb.phase {
			case 0:
				if value := qb.durationList.GetSelection(); value != nil {
					qb.payload.Duration = value.Value.(int)
					qb.phase = 1
					return qb, nil // Ignore tea.Quit from list
				}
			case 1:
				if value := qb.nbPeopleList.GetSelection(); value != nil {
					qb.payload.NbPeople = value.Value.(int)
					qb.phase = 2
					return qb, qb.submitBooking()
				}
			}
		}
	}

	return qb, cmd
}

func (qb *QuickBookModel) View() string {
	switch qb.phase {
	case 0:
		return qb.durationList.View()
	case 1:
		return qb.nbPeopleList.View()
	default:
		return "Booking"
	}
}

func getClosestQuarterHour() time.Time {
	now := time.Now()
	currentHour := now.Hour()
	currentMinutes := now.Minute()

	if currentMinutes > 52 {
		currentHour++
	}

	m1 := math.Round(float64(currentMinutes)/float64(15)) * 15
	m2 := int(m1) % 60

	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		currentHour,
		m2,
		0,
		0,
		now.Location(),
	)
}
func (qb *QuickBookModel) submitBooking() tea.Cmd {
	fmt.Println(qb.payload.Duration, qb.payload.NbPeople)

	return nil
}

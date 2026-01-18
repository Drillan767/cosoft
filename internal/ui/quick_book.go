package ui

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/services"
	"cosoft-cli/internal/ui/components"
	"cosoft-cli/shared/models"
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type QuickBookModel struct {
	// 0 select duration, 1 select nb people, 2 display future loader, 3 display results
	phase int
	// 0 fetch available rooms 1 booking complete
	bookPhase    int
	spinner      spinner.Model
	progress     progress.Model
	durationList components.ListModel
	nbPeopleList components.ListModel
	payload      *api.QuickBookPayload
	rooms        []models.Room
	err          error
}

type bookingProcessStart struct{}
type roomFetchedMsg struct {
	availableRooms []models.Room
}
type noRoomAvailableMsg struct{}
type bookingCompleteMsg struct{}
type bookingFailedMsg struct {
	err error
}

/*
 Rough logic flow:
  phase 1 (enter pressed)
    → return fetchRoomsCmd()
    → set phase = 2 (show spinner + 0% bar)

  Update receives roomsFetchedMsg
    → store rooms in model
    → return bookRoomCmd()
    → set phase = 3 (show spinner + 50% bar)

  Update receives bookingCompleteMsg
    → store result in model
    → set phase = 4 (show results + 100% bar)

  Update receives any error msg
    → store error in model
    → set phase = 5 (error state)

*/

func NewQuickBookModel() *QuickBookModel {
	selection := &api.QuickBookPayload{}
	t := getClosestQuarterHour()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	progress := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(100),
		progress.WithoutPercentage(),
	)

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
		phase:        0,
		durationList: components.NewListModel(durations, "For how long?"),
		nbPeopleList: components.NewListModel(peoples, "For how many people?"),
		payload:      selection,
		spinner:      s,
		progress:     progress,
		rooms:        []models.Room{},
	}
}

func (qb *QuickBookModel) Init() tea.Cmd {
	return nil
}

func (qb *QuickBookModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch qb.phase {
	case 0:
		qb.durationList, cmd = qb.durationList.Update(msg)
	case 1:
		qb.nbPeopleList, cmd = qb.nbPeopleList.Update(msg)
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "enter" {
			switch qb.phase {
			case 0:
				if value := qb.durationList.GetSelection(); value != nil {
					qb.payload.Duration = value.Value.(int)
					qb.payload.DateTime = getClosestQuarterHour()
					qb.phase = 1
					return qb, nil
				}
			case 1:
				if value := qb.nbPeopleList.GetSelection(); value != nil {
					qb.payload.NbPeople = value.Value.(int)
					qb.phase = 2
					qb.bookPhase = 1
					return qb, tea.Batch(qb.spinner.Tick, qb.getRoomsAvailability())
				}
			}
		}

	case roomFetchedMsg:
		qb.bookPhase = 2
		progressCmd := qb.progress.IncrPercent(0.5)
		return qb, tea.Batch(progressCmd, qb.bookRoom())

	case spinner.TickMsg:
		var cmd tea.Cmd
		qb.spinner, cmd = qb.spinner.Update(msg)
		return qb, cmd

	case progress.FrameMsg:
		newModel, cmd := qb.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			qb.progress = newModel
		}

		return qb, cmd
	}

	return qb, cmd
}

func (qb *QuickBookModel) View() string {
	switch qb.phase {
	case 0:
		return qb.durationList.View()
	case 1:
		return qb.nbPeopleList.View()
	case 2:
		spin := qb.spinner.View()

		switch qb.bookPhase {
		case 1:
			spin += " Looking for available rooms... \n\n"
		case 2:
			spin += " Found a meeting room. Booking now... \n\n"
		}

		progress := qb.progress.View()

		err := ""

		if qb.err != nil {
			err = "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(qb.err.Error())
		}

		return spin + progress + err

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

func (qb *QuickBookModel) getRoomsAvailability() tea.Cmd {
	return func() tea.Msg {

		authService, err := services.NewService()

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		user, err := authService.GetAuthData()

		apiClient := api.NewApi()

		payload := api.QuickBookPayload{
			DateTime: qb.payload.DateTime,
			NbPeople: qb.payload.NbPeople,
			Duration: qb.payload.Duration,
		}

		rooms, err := apiClient.GetAvailableRooms(user.WAuth, user.WAuthRefresh, payload)

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		if len(rooms) == 0 {
			return bookingFailedMsg{err: fmt.Errorf("No room available for the selected time")}
		}

		return roomFetchedMsg{availableRooms: rooms}
	}
}

func (qb *QuickBookModel) bookRoom() tea.Cmd {
	return nil
}

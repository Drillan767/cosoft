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
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type QuickBookModel struct {
	// 0 form, 1 display loader, 2 display results
	phase int
	// 0 fetch available rooms, 1 booking in progress, 2 booking complete
	bookPhase  int
	spinner    spinner.Model
	progress   progress.Model
	form       *huh.Form
	payload    *api.CosoftAvailabilityPayload
	rooms      []models.Room
	bookedRoom *models.Room
	err        error
}

type bookingProcessStart struct{}
type roomFetchedMsg struct {
	availableRooms []models.Room
}
type bookingCompleteMsg struct {
	room models.Room
}
type bookingFailedMsg struct {
	err error
}

func NewQuickBookModel() *QuickBookModel {
	selection := &api.CosoftAvailabilityPayload{}
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(100),
		progress.WithoutPercentage(),
	)

	qb := &QuickBookModel{
		phase:    0,
		payload:  selection,
		spinner:  s,
		progress: p,
		rooms:    []models.Room{},
	}

	qb.buildForm()

	return qb
}

func (qb *QuickBookModel) Init() tea.Cmd {
	return qb.form.Init()
}

func (qb *QuickBookModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case roomFetchedMsg:
		qb.rooms = msg.availableRooms
		qb.bookPhase = 2
		progressCmd := qb.progress.IncrPercent(0.5)
		return qb, tea.Batch(progressCmd, qb.bookRoom())

	case bookingCompleteMsg:
		qb.bookedRoom = &msg.room
		qb.bookPhase = 3
		qb.progress = progress.New(
			progress.WithSolidFill("#04B575"),
			progress.WithWidth(100),
			progress.WithoutPercentage(),
		)
		progressCmd := qb.progress.SetPercent(1.0)

		return qb, progressCmd

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

	// Update form during phase 0
	if qb.phase == 0 {
		form, cmd := qb.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			qb.form = f
		}

		if qb.form.State == huh.StateCompleted {
			qb.payload.DateTime = getClosestQuarterHour()
			qb.phase = 1
			qb.bookPhase = 1
			return qb, tea.Batch(qb.spinner.Tick, qb.getRoomsAvailability())
		}

		return qb, cmd
	}

	return qb, nil
}

func (qb *QuickBookModel) View() string {
	switch qb.phase {
	case 0:
		return qb.form.View()
	case 1:
		var header string

		switch qb.bookPhase {
		case 1:
			header = qb.spinner.View() + " Looking for available rooms... \n\n"
		case 2:
			header = qb.spinner.View() + " Found a meeting room. Booking now... \n\n"
		case 3:
			header = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓ Booking complete!") + "\n\n"
		}

		p := qb.progress.View()

		var t string
		var toolTip string
		if qb.bookPhase == 3 && qb.bookedRoom != nil {
			t = qb.generateTable()
			toolTip = "You can now press \"ESC\" to go back to the main menu."
		}

		errMsg := ""
		if qb.err != nil {
			errMsg = "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(qb.err.Error())
		}

		return header + p + t + toolTip + errMsg

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

func (qb *QuickBookModel) buildForm() {
	t := getClosestQuarterHour()

	durations := []components.Item[int]{
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

	peoples := []components.Item[int]{
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

	qb.form = huh.NewForm(
		huh.NewGroup(
			components.NewListField(durations, "For how long?").
				Value(&qb.payload.Duration),
			components.NewListField(peoples, "For how many people?").
				Value(&qb.payload.NbPeople),
		),
	).WithLayout(huh.LayoutStack)
}

func (qb *QuickBookModel) getRoomsAvailability() tea.Cmd {
	return func() tea.Msg {

		authService, err := services.NewService()

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		user, err := authService.GetAuthData()

		apiClient := api.NewApi()

		payload := api.CosoftAvailabilityPayload{
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
	return func() tea.Msg {
		authService, err := services.NewService()

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		user, err := authService.GetAuthData()

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		savedCredits := float64(user.Credits) / float64(100)

		var pickedRoom *models.Room

		for _, room := range qb.rooms {
			if room.NbUsers >= qb.payload.NbPeople {
				pickedRoom = &room
				break
			}
		}

		if pickedRoom == nil {
			return bookingFailedMsg{err: fmt.Errorf("No room suiting user's selection, aborting")}
		}

		if savedCredits < pickedRoom.Price {
			return bookingFailedMsg{err: fmt.Errorf("Not enough credits to perfor, the booking, aborting")}
		}

		payload := api.CosoftBookingPayload{
			CosoftAvailabilityPayload: api.CosoftAvailabilityPayload{
				DateTime: qb.payload.DateTime,
				NbPeople: qb.payload.NbPeople,
				Duration: qb.payload.Duration,
			},
			UserCredits: savedCredits,
			Room:        *pickedRoom,
		}

		apiClient := api.NewApi()

		err = apiClient.BookRoom(user.WAuth, user.WAuthRefresh, payload)

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		return bookingCompleteMsg{room: *pickedRoom}
	}
}

func (qb *QuickBookModel) generateTable() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#fd4b4b"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#fd4b4b")).Bold(true).Align(lipgloss.Center)
			case col == 1:
				return lipgloss.NewStyle().Padding(0, 1).Width(20).Foreground(lipgloss.Color("245"))
			default:
				return lipgloss.NewStyle().Padding(0, 1).Width(14).Foreground(lipgloss.Color("245"))
			}
		}).
		Headers("ROOM", "DURATION", "COST")

	startTime := qb.payload.DateTime
	endTime := startTime.Add(time.Duration(qb.payload.Duration) * time.Minute)
	dateFormat := "02/01/2006 15:04"

	t.Row(
		qb.bookedRoom.Name,
		fmt.Sprintf("%s → %s", startTime.Format(dateFormat), endTime.Format(dateFormat)),
		fmt.Sprintf("%.2f credits", qb.bookedRoom.Price),
	)

	return "\n\n" + t.String() + "\n\n"
}

package ui

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/services"
	"time"

	"cosoft-cli/shared/models"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type LandingModel struct {
	form             *huh.Form
	selection        *models.Selection
	spinner          spinner.Model
	calendarSpinner  spinner.Model
	calendar         string
	nbFutureBookings int
	loading          bool
	loadingCalendar  bool
	err              error
}

type futureBookingMsg struct {
	bookings *api.FutureBookingsResponse
	err      error
}

type updatedCreditsMsg struct {
	credits float64
}

type calendarMsg struct {
	calendar string
	err      error
}

type startFetchingMsg struct{}

func NewLandingModel() *LandingModel {
	selection := &models.Selection{}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	cs := spinner.New()
	cs.Spinner = spinner.Dot
	cs.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	m := &LandingModel{
		selection:        selection,
		spinner:          s,
		loading:          true,
		nbFutureBookings: 0,
		calendar:         "",
		calendarSpinner:  cs,
		loadingCalendar:  true,
	}

	m.buildForm()

	return m
}

func (m *LandingModel) buildForm() {
	resaLabel := "My reservations"
	if m.nbFutureBookings > 0 {
		resaLabel = fmt.Sprintf("My reservations (%d)", m.nbFutureBookings)
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What's the plan?").
				Options(
					huh.NewOption("Quick book", "quick-book"),
					huh.NewOption("Browse & book", "browse"),
					huh.NewOption(resaLabel, "reservations"),
					huh.NewOption("Previous reservations", "history"),
					huh.NewOption("Settings", "settings"),
					huh.NewOption("Quit", "quit"),
				).
				Value(&m.selection.Choice),
		),
	)
}

func (m *LandingModel) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		m.spinner.Tick,
		m.calendarSpinner.Tick,
		m.form.Init(),
		m.spinner.Tick,
		m.calendarSpinner.Tick,
		func() tea.Msg {
			return startFetchingMsg{}
		},
	)
}

func (m *LandingModel) fetchFutureBookings() tea.Cmd {
	return func() tea.Msg {
		authService, err := services.NewService()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		user, err := authService.GetAuthData()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		apiClient := api.NewApi()
		b, err := apiClient.GetFutureBookings(user.WAuth, user.WAuthRefresh)

		return futureBookingMsg{bookings: b, err: err}
	}
}

func (m *LandingModel) updateCredits() tea.Cmd {
	return func() tea.Msg {
		s, err := services.NewService()

		if err != nil {
			return nil
		}

		credits, err := s.UpdateCredits()

		// Either failed, or there was nothing to update
		if err != nil || credits == nil {
			return nil
		}

		return updatedCreditsMsg{credits: *credits}
	}
}

func (m *LandingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startFetchingMsg:
		return m, tea.Batch(
			m.fetchFutureBookings(),
			m.getCalendarView(),
			m.updateCredits(),
		)

	case futureBookingMsg:
		if msg.err != nil {
			m.loading = false
			m.buildForm()
			return m, m.form.Init()
		}
		m.nbFutureBookings = msg.bookings.Total
		m.loading = false
		m.buildForm() // Rebuild form with updated data
		return m, m.form.Init()

	case calendarMsg:
		m.loadingCalendar = false

		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.calendar = msg.calendar

	case updatedCreditsMsg:
		credits := fmt.Sprintf("Credits: %.02f", msg.credits)
		return m, func() tea.Msg {
			return UpdateHeaderMsg{Credits: &credits}
		}

	case spinner.TickMsg:
		var cmd1, cmd2 tea.Cmd
		m.spinner, cmd1 = m.spinner.Update(msg)
		m.calendarSpinner, cmd2 = m.calendarSpinner.Update(msg)
		return m, tea.Batch(cmd1, cmd2)
	}

	if m.loading {
		return m, nil
	}

	form, cmd := m.form.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		return m, func() tea.Msg {
			return NavigateMsg{Page: m.selection.Choice}
		}
	}

	return m, cmd
}

func (m *LandingModel) View() string {
	// Calendar
	var calendar string
	var loadingCalendar string
	// var loadingMenu string

	if m.loadingCalendar {
		loadingCalendar = fmt.Sprintf("%s Loading calendar informations...\n\n", m.calendarSpinner.View())
	}

	if m.calendar != "" {
		calendar = m.calendar + "\n\n"
	}

	if m.loading {
		return loadingCalendar + fmt.Sprintf("\n %s Loading...\n", m.spinner.View())
	}
	if m.form == nil {
		return "Error: form is nil"
	}
	return loadingCalendar + calendar + m.form.View()
}

func (m *LandingModel) GetSelection() *models.Selection {
	return m.selection
}

func (m *LandingModel) getCalendarView() tea.Cmd {
	return func() tea.Msg {
		authService, err := services.NewService()

		if err != nil {
			return calendarMsg{err: err}
		}

		err = authService.EnsureRoomsStored()

		if err != nil {
			return calendarMsg{err: err}
		}

		now := time.Now()
		date := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0,
			0,
			0,
			0,
			now.Location(),
		)

		usage, err := authService.GetRoomAvailabilities(date)

		if err != nil {
			return calendarMsg{err: err}
		}

		calendar := ""

		for _, u := range usage {
			calendar = fmt.Sprintf("%s\n %s: ", m.calendar, u.Name)
		}

		return calendarMsg{calendar: calendar}
	}
}

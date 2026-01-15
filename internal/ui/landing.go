package ui

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/auth"
	"cosoft-cli/shared/models"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type LandingModel struct {
	form             *huh.Form
	selection        *models.Selection
	spinner          spinner.Model
	nbFutureBookings int
	loading          bool
}

type futureBookingMsg struct {
	total int
	err   error
}

func NewLandingModel() *LandingModel {
	selection := &models.Selection{}

	s := spinner.New()
	s.Spinner = spinner.Dot

	m := &LandingModel{
		selection:        selection,
		spinner:          s,
		loading:          true,
		nbFutureBookings: 0,
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
					huh.NewOption(resaLabel, "resa"),
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
		m.fetchFutureBookings(),
	)
}

func (m *LandingModel) fetchFutureBookings() tea.Cmd {
	return func() tea.Msg {
		authService, err := auth.NewAuthService()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		user, err := authService.GetAuthData()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		apiClient := api.NewApi()

		total, err := apiClient.GetFutureBookings(user.WAuth, user.WAuthRefresh)

		return futureBookingMsg{total: total, err: err}
	}
}

func (m *LandingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case futureBookingMsg:
		if msg.err != nil {
			m.loading = false
			m.buildForm()
			return m, m.form.Init()
		}
		m.nbFutureBookings = msg.total
		m.loading = false
		m.buildForm() // Rebuild form with updated data
		return m, m.form.Init()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
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
	if m.loading {
		return fmt.Sprintf("\n %s Loading...\n", m.spinner.View())
	}
	if m.form == nil {
		return "Error: form is nil"
	}
	return m.form.View()
}

func (m *LandingModel) GetSelection() *models.Selection {
	return m.selection
}

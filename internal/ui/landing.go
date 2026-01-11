package ui

import (
	"cosoft-cli/shared/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type LandingModel struct {
	form      *huh.Form
	selection *models.Selection
}

func NewLandingModel() *LandingModel {
	selection := &models.Selection{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What's the plan?").
				Options(
					huh.NewOption("Quick book", "quick-book"),
					huh.NewOption("Browse & book", "browse"),
					huh.NewOption("My reservations", "resa"),
					huh.NewOption("Previous reservations", "history"),
					huh.NewOption("Settings", "settings"),
					huh.NewOption("Quit", "quit"),
				).
				Value(&selection.Choice),
		),
	)

	return &LandingModel{
		form:      form,
		selection: selection,
	}
}

func (m *LandingModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *LandingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return m.form.View()
}

func (m *LandingModel) GetSelection() *models.Selection {
	return m.selection
}

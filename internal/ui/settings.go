package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type SettingsModel struct {
	spinner   spinner.Model
	form      *huh.Form
	confirmed bool
	choice    string
	loading   bool
	err       error
}

func NewSettingsModel() *SettingsModel {

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	settings := &SettingsModel{
		choice:    "",
		confirmed: false,
		spinner:   s,
		loading:   false,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a setting").
				Options(
					huh.NewOption("Delete the local settings", "clean"),
				).
				Value(&settings.choice),
			huh.NewNote().
				Title("⚠️  Warning").
				Description(
					"You are about to delete the local settings file, which will also log you out.\n"+
						"This cannot be undone.",
				),
			huh.NewConfirm().
				Title("Confirm?").
				Negative("No").
				Affirmative("Yes").
				Value(&settings.confirmed),
		),
	)

	settings.form = form

	return settings
}

func (m *SettingsModel) Init() tea.Cmd {
	return m.form.Init()
}

func (s *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}

	form, cmd := s.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		s.form = f
	}

	return s, cmd
}

func (s *SettingsModel) View() string {
	return s.form.View()
}

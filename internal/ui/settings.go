package ui

import (
	"cosoft-cli/internal/services"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type SettingsModel struct {
	phase       int
	spinner     spinner.Model
	choiceForm  *huh.Form
	confirmForm *huh.Form
	confirmed   bool
	choice      string
	loading     bool
	err         error
}

type clearingDone struct {
	err error
}

func NewSettingsModel() *SettingsModel {

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	settings := &SettingsModel{
		phase:     1,
		choice:    "",
		confirmed: false,
		spinner:   s,
		loading:   false,
	}

	choice := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a setting").
				Options(
					huh.NewOption("Delete the local settings", "clean"),
				).
				Value(&settings.choice),
		),
	)

	confirm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm?").
				Negative("No").
				Affirmative("Yes").
				Value(&settings.confirmed),
		),
	)

	settings.choiceForm = choice
	settings.confirmForm = confirm

	return settings
}

func (s *SettingsModel) Init() tea.Cmd {
	return s.choiceForm.Init()
}

func (s *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	case clearingDone:
		s.loading = false
		if msg.err != nil {
			s.err = msg.err
			return s, nil
		}

		s.phase = 4
		return s, tea.Printf("")
	}

	switch s.phase {
	case 1:
		form, cmd := s.choiceForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			s.choiceForm = f
		}
		if s.choiceForm.State == huh.StateCompleted {
			s.phase = 2
			return s, s.confirmForm.Init()
		}
		return s, cmd
	case 2:
		form, cmd := s.confirmForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			s.confirmForm = f
		}

		if s.confirmForm.State == huh.StateCompleted {
			s.phase = 3
			return s, tea.Batch(
				s.spinner.Tick,
				s.clearInformation(),
			)
		}
		return s, cmd
	}

	return s, cmd
}

func (s *SettingsModel) View() string {
	if s.err != nil {
		return s.err.Error()
	}

	switch s.phase {
	case 1:
		return s.choiceForm.View()
	case 2:
		var w string
		if s.choice == "clean" {
			w = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")).
				Bold(true).
				Render("⚠️  Warning") +
				"\n\n" +
				"You are about to delete the local settings file, which will also log you out. \n" +
				"This cannot be undone." +
				"\n\n"
		}

		return w + s.confirmForm.View()
	case 3:
		return s.spinner.View() + " Clearing personal data..."
	case 4:
		success := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render("✓ Cleared Personal Data!")

		tooltip := "You can now press \"ESC\" to quit the program."

		return success + "\n\n" + tooltip
	default:
		return "Settings"
	}
}

func (s *SettingsModel) ShouldQuitOnEsc() bool {
	return s.phase == 4
}

func (s *SettingsModel) clearInformation() tea.Cmd {
	return func() tea.Msg {
		clearService, err := services.NewService()

		if err != nil {
			return clearingDone{err: err}
		}

		err = clearService.ClearData()
		return clearingDone{err: err}
	}
}

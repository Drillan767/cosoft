package ui

import (
	"cosoft-cli/internal/api"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func required(s string) error {
	if s == "" {
		return errors.New("Field is required")
	}

	return nil
}

// LoginModel is a Bubbletea model for the login form
type LoginModel struct {
	form        *huh.Form
	credentials *api.LoginPayload
	quitting    bool
}

// NewLoginModel creates a new login form model
func NewLoginModel() *LoginModel {
	creds := &api.LoginPayload{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Login").
				Description("Please insert your credentials to authenticate"),
			huh.NewInput().
				Validate(required).
				Title("Email address").
				Value(&creds.Email),
			huh.NewInput().
				EchoMode(huh.EchoModePassword).
				Validate(required).
				Title("Password").
				Value(&creds.Password),
		),
	)

	return &LoginModel{
		form:        form,
		credentials: creds,
	}
}

func (m *LoginModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Update form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check if form is complete
	if m.form.State == huh.StateCompleted {
		m.quitting = true
		return m, tea.Quit
	}

	return m, cmd
}

func (m *LoginModel) View() string {
	if m.quitting {
		return ""
	}
	return m.form.View()
}

// GetCredentials returns the entered credentials
func (m *LoginModel) GetCredentials() *api.LoginPayload {
	return m.credentials
}

// LoginFormWithLayout creates a login form wrapped in a layout
func (ui *UI) LoginFormWithLayout() (*api.LoginPayload, error) {
	loginModel := NewLoginModel()

	// Wrap in layout
	layout := NewLayoutWithDefaults(
		loginModel,
		"COSOFT CLI - Authentication",
		"Press Ctrl+C to cancel",
	)

	p := tea.NewProgram(layout)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	// Extract the login model from the layout
	if layoutModel, ok := finalModel.(*Layout); ok {
		if loginModel, ok := layoutModel.GetContent().(*LoginModel); ok {
			return loginModel.GetCredentials(), nil
		}
	}

	return nil, fmt.Errorf("failed to retrieve login credentials")
}

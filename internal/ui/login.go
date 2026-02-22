package ui

import (
	"cosoft-cli/internal/api"
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// LoginModel is a Bubbletea model for the login form
type LoginModel struct {
	form        *huh.Form
	credentials *api.LoginPayload
	user        *api.UserResponse
	spinner     spinner.Model
	err         error
	quitting    bool
	loading     bool
}

type loginSuccessMsg struct {
	user *api.UserResponse
}

type loginErrorMsg struct {
	err error
}

func (m *LoginModel) Init() tea.Cmd {
	return m.form.Init()
}

func required(s string) error {
	if s == "" {
		return errors.New("field is required")
	}

	return nil
}

func (m *LoginModel) performLogin() tea.Cmd {
	return func() tea.Msg {
		apiClient := api.NewApi()

		user, err := apiClient.Login(m.credentials)

		if err != nil {
			return loginErrorMsg{err: err}
		}

		return loginSuccessMsg{user: user}
	}
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

	s := spinner.New()
	s.Spinner = spinner.Dot

	return &LoginModel{
		form:        form,
		credentials: creds,
		spinner:     s,
		loading:     false,
	}
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case loginSuccessMsg:
		m.user = msg.user // Store the user!
		m.loading = false
		m.quitting = true
		return m, tea.Quit

	case loginErrorMsg:
		m.loading = false
		m.err = msg.err
		// Reset form to allow retry
		m.form = huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("Login").
					Description("Please insert your credentials to authenticate"),
				huh.NewInput().
					Validate(required).
					Title("Email address").
					Value(&m.credentials.Email),
				huh.NewInput().
					EchoMode(huh.EchoModePassword).
					Validate(required).
					Title("Password").
					Value(&m.credentials.Password),
			),
		)
		return m, m.form.Init()
	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if m.loading {
		return m, nil
	}

	// Update form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check if form is complete
	if m.form.State == huh.StateCompleted {
		m.loading = true
		return m, tea.Batch(
			m.spinner.Tick,
			m.performLogin(),
		)
	}

	return m, cmd
}

func (m *LoginModel) View() string {
	if m.quitting {
		return ""
	}

	if m.loading {
		return fmt.Sprintf("\n%s Logging in...", m.spinner.View())
	}

	if m.err != nil {
		return m.form.View() + fmt.Sprintf("\n\n❌ Error: %v\n", m.err)
	}

	return m.form.View()
}

// GetCredentials returns the entered credentials
func (m *LoginModel) GetCredentials() *api.LoginPayload {
	return m.credentials
}

func (m *LoginModel) GetUser() *api.UserResponse {
	return m.user
}

// LoginForm creates a login form wrapped in a layout
func (ui *UI) LoginForm() (*LoginModel, error) {
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
			return loginModel, nil
		}
	}

	return nil, fmt.Errorf("failed to retrieve login credentials")
}

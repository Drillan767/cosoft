package ui

import (
	"cosoft-cli/internal/auth"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) StartApp(startPage string, allowBackNav bool) error {
	appModel := NewAppModel(startPage, allowBackNav)

	header := "COSOFT CLI"

	// Try to get user info
	authService, err := auth.NewAuthService()

	if err != nil {
		return err
	}

	if user, err := authService.GetAuthData(); err == nil {
		header = fmt.Sprintf("COSOFT CLI | %s %s (%s) | Credits: %d",
			user.FirstName, user.LastName, user.Email, user.Credits/100)
	}

	layout := NewLayoutWithDefaults(
		appModel,
		header,
		"Press Ctrl + C to cancel",
	)

	p := tea.NewProgram(layout)

	_, err = p.Run()

	return err
}

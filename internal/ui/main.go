package ui

import (
	"cosoft-cli/internal/services"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) StartApp(startPage string, allowBackNav bool) error {
	appModel := NewAppModel(startPage, allowBackNav)

	config := DefaultLayoutConfig()
	config.Header.Left = "COSOFT CLI"
	config.Header.Center = startPage // Initial location
	config.Footer = "Press Ctrl + C to cancel"

	// Try to get user info
	authService, err := services.NewService()

	if err != nil {
		return err
	}

	if user, err := authService.GetAuthData(); err == nil {
		credits := float32(user.Credits) / float32(100)
		config.Header.Right = fmt.Sprintf("%s %s (%s)", user.FirstName, user.LastName, user.Email)
		config.Header.Credits = fmt.Sprintf("Credits: %.02f", credits)
	}

	layout := NewLayout(appModel, config)

	p := tea.NewProgram(layout)

	_, err = p.Run()

	return err
}

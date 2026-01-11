package ui

import tea "github.com/charmbracelet/bubbletea"

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) StartApp(startPage string, allowBackNav bool) error {
	appModel := NewAppModel(startPage, allowBackNav)

	layout := NewLayoutWithDefaults(
		appModel,
		"COSOFT CLI",
		"Press Ctrl + C to cancel",
	)

	p := tea.NewProgram(layout)

	_, err := p.Run()

	return err
}

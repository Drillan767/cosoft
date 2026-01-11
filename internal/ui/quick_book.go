package ui

import (
	"cosoft-cli/internal/api"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type QuickBookModel struct {
	form     *huh.Form
	duration *api.QuickBookPayload
}

func NewQuickBookModel() *QuickBookModel {
	selection := &api.QuickBookPayload{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Pick a duration").
				Options(
					huh.NewOption("30 minutes", 30),
					huh.NewOption("1 hour", 60),
				).
				Value(&selection.Duration),
		),
	)

	return &QuickBookModel{
		form:     form,
		duration: selection,
	}
}

func (qb *QuickBookModel) Init() tea.Cmd {
	return qb.form.Init()
}

func (qb *QuickBookModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := qb.form.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		qb.form = f
	}

	if qb.form.State == huh.StateCompleted {
		// Do the thing.
	}

	return qb, cmd
}

func (qb *QuickBookModel) View() string {
	if qb.form.State == huh.StateCompleted {
		return fmt.Sprintf("Selected: %s", qb.form.GetString("duration"))
	}
	return qb.form.View()
}

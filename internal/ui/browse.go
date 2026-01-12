package ui

import (
	"cosoft-cli/internal/api"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type BrowseModel struct {
	form    *huh.Form
	choices *api.BrowsePayload
}

func NewBrowseModel() *BrowseModel {
	choices := &api.BrowsePayload{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Reservation date").
				Description("Pick a date in the future, format yyyy-mm-dd").
				Validate(validateDateIsFuture).
				Value(&choices.StarDate),
			huh.NewInput().
				Title("Reservation hour").
				Description("The hour needs to be rounded to the quarter (ex: 9:15, 10:30, etc)").
				Validate(validateHour).
				Value(&choices.StartHour),
		),
	)

	return &BrowseModel{
		form:    form,
		choices: choices,
	}
}

func validateDateIsFuture(s string) error {
	if s == "" {
		return fmt.Errorf("Date is required")
	}

	date, err := time.Parse(time.DateOnly, s)

	if err != nil {
		return fmt.Errorf("Date could not be parsed")
	}

	if time.Now().After(date) {
		return fmt.Errorf("Date is not in the future")
	}

	return nil
}

func validateHour(s string) error {
	if s == "" {
		return fmt.Errorf("Hour is required")
	}

	time, err := time.Parse("15:04", s)

	if err != nil {
		return fmt.Errorf("Could not parse time")
	}

	hours := time.Hour()

	if hours < 8 || hours > 20 {
		return fmt.Errorf("Hours outside opening hours")
	}

	minutes := time.Minute()

	if minutes%15 != 0 {
		return fmt.Errorf("Minutes not rounded to quarters")
	}

	return nil
}

func (b *BrowseModel) Init() tea.Cmd {
	return b.form.Init()
}

func (b *BrowseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := b.form.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		b.form = f
	}

	if b.form.State == huh.StateCompleted {
		// Step 2 or whatever
	}

	return b, cmd
}

func (b *BrowseModel) View() string {
	return b.form.View()
}

func (b *BrowseModel) GetChoice() *api.BrowsePayload {
	return b.choices
}

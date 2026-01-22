package ui

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/shared/models"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type BrowseModel struct {
	phase      int
	spinner    spinner.Model
	rooms      []models.Room
	bookedRoom *models.Room
	searchForm *huh.Form
	bookForm   *huh.Form
	choices    *api.BrowsePayload
}

/*
Step 1
- Get
*/

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
			huh.NewSelect[int]().
				Title("Reservation duration").
				Options(
					huh.NewOption("30mn", 30),
					huh.NewOption("1 hour", 60),
				).
				Value(&choices.Duration),
		),
	)

	return &BrowseModel{
		searchForm: form,
		choices:    choices,
	}
}

func validateDateIsFuture(s string) error {
	if s == "" {
		return fmt.Errorf("date is required")
	}

	date, err := time.Parse(time.DateOnly, s)

	if err != nil {
		return fmt.Errorf("date could not be parsed")
	}

	if time.Now().After(date) {
		return fmt.Errorf("date is not in the future")
	}

	return nil
}

func validateHour(s string) error {
	if s == "" {
		return fmt.Errorf("hour is required")
	}

	h, err := time.Parse("15:04", s)

	if err != nil {
		return fmt.Errorf("could not parse time")
	}

	hours := h.Hour()

	if hours < 8 || hours > 20 {
		return fmt.Errorf("hours outside opening hours")
	}

	minutes := h.Minute()

	if minutes%15 != 0 {
		return fmt.Errorf("minutes not rounded to quarters")
	}

	return nil
}

func (b *BrowseModel) Init() tea.Cmd {
	return b.searchForm.Init()
}

func (b *BrowseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := b.searchForm.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		b.searchForm = f
	}

	if b.searchForm.State == huh.StateCompleted {
		// Step 2 or whatever
	}

	return b, cmd
}

func (b *BrowseModel) View() string {
	return b.searchForm.View()
}

func (b *BrowseModel) GetChoice() *api.BrowsePayload {
	return b.choices
}

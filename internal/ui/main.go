package ui

import (
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/services"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const timeOnlyFormat = "15:04"

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (ui *UI) StartApp(startPage string, allowBackNav bool) error {
	appModel := NewAppModel(startPage, allowBackNav)

	config := DefaultLayoutConfig()
	config.Header.Left = "COSOFT CLI"
	config.Header.Center = strings.ToUpper(startPage[:1]) + startPage[1:] // Initial location
	config.Footer = "Press Ctrl + C to cancel"

	// Try to get user info
	authService, err := services.NewService()
	if err != nil {
		return err
	}

	if user, err := authService.GetAuthData(); err == nil {
		config.Header.Right = fmt.Sprintf("%s %s (%s)", user.FirstName, user.LastName, user.Email)
		config.Header.Credits = fmt.Sprintf("Credits: %.02f", user.Credits)
	}

	layout := NewLayout(appModel, config)

	p := tea.NewProgram(layout)

	_, err = p.Run()

	return err
}

func validateDateIsFuture(s string) error {
	if s == "" {
		return fmt.Errorf("date is required")
	}

	location, err := common.LoadLocalTime()
	if err != nil {
		return err
	}

	date, err := time.ParseInLocation(time.DateOnly, s, location)
	if err != nil {
		return fmt.Errorf("date could not be parsed")
	}

	// Compute the today's date and compare it with the input.
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if date.Before(today) {
		return fmt.Errorf("date is not in the future")
	}

	return nil
}

func validateHour(s string) error {
	if s == "" {
		return fmt.Errorf("hour is required")
	}

	h, err := time.Parse(timeOnlyFormat, s)
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

// roundHourToQuarter rounds t to the next quarter hour.
func roundHourToQuarter(t time.Time) time.Time {
	return t.Truncate(15 * time.Minute).Add(15 * time.Minute)
}

package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Reservation struct {
	OrderResourceRentId string
	ItemName            string
	Start               string
	End                 string
	Credits             float64
}

type ReservationListModel struct {
	phase        int
	spinner      spinner.Model
	bookingId    string
	reservations []Reservation
	err          error
}

func NewReservationListModel() *ReservationListModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return &ReservationListModel{
		phase:   0,
		spinner: s,
	}
}

func (rl *ReservationListModel) Init() tea.Cmd {
	// TODO: Fetch reservations, assign them to model
	return nil
}

func (rl *ReservationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return rl, nil
}

func (rl *ReservationListModel) View() string {
	return "My reservations"
}

package ui

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/services"
	"cosoft-cli/internal/ui/components"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type ReservationListModel struct {
	loading         bool
	phase           int
	spinner         spinner.Model
	bookingId       string
	reservationList *components.ListField[string]
	confirm         *huh.Form
	comfirmed       bool
	reservations    api.FutureBookingsResponse
	err             error
}

type cancelComplete struct{}

func NewReservationListModel() *ReservationListModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	rl := &ReservationListModel{
		phase:     0,
		loading:   true,
		comfirmed: true,
		spinner:   s,
	}

	rl.buildConfirmForm()

	return rl
}

func (rl *ReservationListModel) Init() tea.Cmd {
	return tea.Batch(
		rl.spinner.Tick,
		rl.fetchFutureBookings(),
		rl.confirm.Init(),
	)
}

func (rl *ReservationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch rl.phase {
	case 1:
		m, c := rl.reservationList.Update(msg)
		rl.reservationList = m.(*components.ListField[string])
		cmd = c
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "enter" && rl.phase == 1 {
			if rl.reservationList.IsSelected() {
				value := rl.reservationList.SelectedItem()
				rl.bookingId = value.Value
				rl.phase = 2
				return rl, nil
			}
		}

	case futureBookingMsg:
		rl.loading = false

		if msg.err != nil {
			rl.err = msg.err
			return rl, nil
		}

		rl.reservations = *msg.bookings
		if err := rl.buildReservationList(); err != nil {
			rl.err = err
			return rl, nil
		}

		rl.phase = 1
		return rl, nil

	case cancelComplete:
		rl.loading = false
		rl.phase = 4
		return rl, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		rl.spinner, cmd = rl.spinner.Update(msg)
		return rl, cmd
	}

	form, cmd := rl.confirm.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		rl.confirm = f
	}

	if rl.confirm.State == huh.StateCompleted {
		rl.phase = 3
		rl.loading = true
		return rl, tea.Batch(rl.spinner.Tick, rl.cancelReservation())
	}

	return rl, cmd
}

func (rl *ReservationListModel) View() string {
	if rl.err != nil {
		return rl.err.Error()
	}

	switch rl.phase {
	case 0:
		if rl.loading {
			return fmt.Sprintf("\n %s Loading reservations...\n", rl.spinner.View())
		}
	case 1:
		return rl.reservationList.View()
	case 2:
		return rl.confirm.View()
	case 3:
		if rl.loading {
			return fmt.Sprintf("\n %s Cancelling reservation...\n", rl.spinner.View())
		}
	case 4:
		success := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓ Booking cancelled successfully!")
		toolTip := "You can now press \"ESC\" to go back to the main menu."

		return success + "\n\n" + toolTip
	}

	return "My reservations"
}

func (rl *ReservationListModel) fetchFutureBookings() tea.Cmd {
	return func() tea.Msg {
		authService, err := services.NewService()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		user, err := authService.GetAuthData()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		apiClient := api.NewApi()
		bookings, err := apiClient.GetFutureBookings(user.WAuth, user.WAuthRefresh)

		return futureBookingMsg{
			bookings: bookings,
			err:      err,
		}
	}
}

func (rl *ReservationListModel) cancelReservation() tea.Cmd {
	return func() tea.Msg {
		authService, err := services.NewService()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		user, err := authService.GetAuthData()

		if err != nil {
			return futureBookingMsg{err: err}
		}

		apiClient := api.NewApi()
		err = apiClient.CancelBooking(user.WAuth, user.WAuthRefresh, rl.bookingId)

		if err != nil {
			return futureBookingMsg{err: err}
		}

		return cancelComplete{}
	}
}

func (rl *ReservationListModel) buildReservationList() error {
	list := make([]components.Item[string], len(rl.reservations.Data))

	for i, r := range rl.reservations.Data {
		parsedStart, err := time.Parse("2006-01-02T15:04:05", r.Start)

		if err != nil {
			return err
		}

		parsedEnd, err := time.Parse("2006-01-02T15:04:05", r.End)

		if err != nil {
			return err
		}

		dateFormat := "02/01/2006 15:04"

		list[i] = components.Item[string]{
			Label: r.ItemName,
			Value: r.OrderResourceRentId,
			Subtitle: fmt.Sprintf(
				"%s → %s · %.02f credits",
				parsedStart.Format(dateFormat),
				parsedEnd.Format(dateFormat),
				r.Credits,
			),
		}
	}

	rl.reservationList = components.NewListField(list, "Pick a reservation to cancel it")

	return nil
}

func (rl *ReservationListModel) buildConfirmForm() {
	rl.confirm = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm cancellation?").
				Affirmative("Yes").
				Negative("No").
				Value(&rl.comfirmed),
		),
	)
}

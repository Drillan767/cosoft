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
	phase        int
	confirmed    bool
	bookingId    string
	reservations api.FutureBookingsResponse
	form         *huh.Form
	spinner      spinner.Model
	err          error
}

type cancelComplete struct{}

func NewReservationListModel() *ReservationListModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return &ReservationListModel{
		phase:     1,
		confirmed: false,
		spinner:   s,
		bookingId: "",
	}
}

func (rl *ReservationListModel) Init() tea.Cmd {
	return tea.Batch(
		rl.spinner.Tick,
		rl.fetchFutureBookings(),
	)
}

func (rl *ReservationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		rl.spinner, cmd = rl.spinner.Update(msg)
		return rl, cmd

	case futureBookingMsg:
		if msg.err != nil {
			rl.err = msg.err
			return rl, nil
		}

		rl.reservations = *msg.bookings
		rl.phase = 2

		if err := rl.buildForm(); err != nil {
			rl.err = err
			return rl, nil
		}

		return rl, rl.form.Init()

	case cancelComplete:
		rl.phase = 4
		return rl, nil
	}

	if rl.form == nil {
		return rl, nil
	}

	form, cmd := rl.form.Update(msg)

	if f, ok := form.(*huh.Form); ok {
		rl.form = f
	}

	if rl.form.State == huh.StateCompleted {
		rl.phase = 3
		return rl, tea.Batch(rl.spinner.Tick, rl.cancelReservation())
	}

	return rl, cmd
}

func (rl *ReservationListModel) View() string {
	if rl.err != nil {
		return rl.err.Error()
	}

	switch rl.phase {
	case 1:
		return rl.spinner.View() + " Loading reservations..."
	case 2:
		if len(rl.reservations.Data) == 0 {
			return "No reservations found \n\n Press \"ESC\" to go back to the main menu."
		}
		return rl.form.View()
	case 3:
		return rl.spinner.View() + " Cancelling reservations..."
	case 4:
		success := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render("✓ Cancellation complete!")

		tooltip := "You can now press \"ESC\" to go back to the main menu."

		return success + "\n\n" + tooltip
	default:
		return "Reservations"

	}
}

func (rl *ReservationListModel) buildForm() error {
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

		duration := parsedEnd.Sub(parsedStart).Minutes()

		paidPrice := r.Credits * (float64(duration) / 60)

		dateFormat := "02/01/2006 15:04"

		list[i] = components.Item[string]{
			Label: r.ItemName,
			Value: r.OrderResourceRentId,
			Subtitle: fmt.Sprintf(
				"%s → %s · %.02f credits",
				parsedStart.Format(dateFormat),
				parsedEnd.Format(dateFormat),
				paidPrice,
			),
		}
	}

	rl.form = huh.NewForm(
		huh.NewGroup(
			components.NewListField(list, "Pick a reservation to cancel it").
				Value(&rl.bookingId),
			huh.NewConfirm().
				Title("Confirm cancellation?").
				Negative("No").
				Affirmative("Yes").
				Value(&rl.confirmed),
		),
	)

	return nil
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

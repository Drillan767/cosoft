package ui

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/services"
	"cosoft-cli/internal/ui/components"
	"cosoft-cli/shared/models"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type BrowseModel struct {
	phase         int
	spinner       spinner.Model
	rooms         []models.Room
	roomId        string
	bookedRoom    *models.Room
	searchForm    *huh.Form
	bookForm      *huh.Form
	browsePayload *api.BrowsePayload
	bookPayload   *api.CosoftBookingPayload
	err           error
}

func NewBrowseModel() *BrowseModel {
	browsePayload := &api.BrowsePayload{
		StartDate: time.Now().Format(time.DateOnly),
		StartHour: roundHourToQuarter(time.Now()).Format(timeOnlyFormat),
	}
	bookPayload := &api.CosoftBookingPayload{}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	peoples := []components.Item[int]{
		{
			Label:    "1 person",
			Subtitle: "The research will include callboxes",
			Value:    1,
		},
		{
			Label:    "2 persons or more",
			Subtitle: "The research will default to classic meeting rooms",
			Value:    2,
		},
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Reservation date").
				Description("Pick a date in the future, format yyyy-mm-dd").
				Validate(validateDateIsFuture).
				Value(&browsePayload.StartDate),
			huh.NewInput().
				Title("Reservation hour").
				Description("The hour needs to be rounded to the quarter (ex: 9:15, 10:30, etc)").
				Validate(validateHour).
				Value(&browsePayload.StartHour),
			huh.NewSelect[int]().
				Title("Reservation duration").
				Options(
					huh.NewOption("30mn", 30),
					huh.NewOption("1 hour", 60),
					huh.NewOption("1 hour 30 minutes", 90),
					huh.NewOption("2 hours", 120),
				).
				Value(&browsePayload.Duration),
			components.NewListField(peoples, "For how many people?").
				Value(&browsePayload.NbPeople),
		),
	)

	return &BrowseModel{
		phase:         0,
		spinner:       s,
		searchForm:    form,
		roomId:        "",
		browsePayload: browsePayload,
		bookPayload:   bookPayload,
	}
}

func (b *BrowseModel) Init() tea.Cmd {
	b.browsePayload.StartDate = time.Now().Format(time.DateOnly)
	b.browsePayload.StartHour = roundHourToQuarter(time.Now()).Format(timeOnlyFormat)
	return b.searchForm.Init()
}

func (b *BrowseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case bookingFailedMsg:
		b.err = msg.err
	case roomFetchedMsg:
		b.rooms = msg.availableRooms
		b.bookForm = b.buildBookForm(b.rooms)
		b.phase = 2
		return b, b.bookForm.Init()
	case bookingCompleteMsg:
		b.bookedRoom = &msg.room
		b.phase = 4
		return b, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		b.spinner, cmd = b.spinner.Update(msg)
		return b, cmd
	}

	switch b.phase {
	case 0:
		form, cmd := b.searchForm.Update(msg)

		if f, ok := form.(*huh.Form); ok {
			b.searchForm = f
		}

		if b.searchForm.State == huh.StateCompleted {
			b.phase = 1
			return b, tea.Batch(b.spinner.Tick, b.getRoomsAvailability())
		}

		return b, cmd
	case 2:
		form, cmd := b.bookForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			b.bookForm = f
		}

		if b.bookForm.State == huh.StateCompleted {
			b.phase = 3
			return b, tea.Batch(b.spinner.Tick, b.bookRoom())
		}

		return b, cmd
	case 3:

	}

	return b, nil
}

func (b *BrowseModel) View() string {

	if b.err != nil {
		return b.err.Error()
	}

	switch b.phase {
	case 0:
		return b.searchForm.View()
	case 1:
		return b.spinner.View() + " Looking for available rooms... \n\n"
	case 2:
		return b.bookForm.View()
	case 3:
		return b.spinner.View() + " Booking selected room... \n\n"
	case 4:
		if b.err != nil {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(b.err.Error())
		}

		header := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓ Booking complete!") + "\n\n"
		tooltip := "You can now press \"ESC\" to go back to the main menu."
		t := b.generateTable()

		return header + tooltip + t

	}

	return "Browse."
}

func (b *BrowseModel) getRoomsAvailability() tea.Cmd {
	return func() tea.Msg {

		authService, err := services.NewService()

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		user, err := authService.GetAuthData()

		apiClient := api.NewApi()

		dt := b.getStartTime(b.browsePayload.StartDate, b.browsePayload.StartHour)

		payload := api.CosoftAvailabilityPayload{
			DateTime: dt,
			NbPeople: b.browsePayload.NbPeople,
			Duration: b.browsePayload.Duration,
		}

		rooms, err := apiClient.GetAvailableRooms(user.WAuth, user.WAuthRefresh, payload)

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		if len(rooms) == 0 {
			return bookingFailedMsg{err: fmt.Errorf("no room available for the selected time")}
		}

		return roomFetchedMsg{availableRooms: rooms}
	}
}

func (b *BrowseModel) GetChoice() *api.BrowsePayload {
	return b.browsePayload
}

func (b *BrowseModel) buildBookForm(rooms []models.Room) *huh.Form {
	list := make([]components.Item[string], len(rooms))

	for i, room := range rooms {
		list[i] = components.Item[string]{
			Value:    room.Id,
			Label:    room.Name,
			Subtitle: fmt.Sprintf("%.02f credits", room.Price),
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			components.NewListField(list, "Pick a meeting room").
				Value(&b.roomId),
		))

	return form
}

func (b *BrowseModel) bookRoom() tea.Cmd {
	return func() tea.Msg {
		authService, err := services.NewService()

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		user, err := authService.GetAuthData()

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		var pickedRoom *models.Room

		for _, room := range b.rooms {
			if room.Id == b.roomId {
				pickedRoom = &room
				break
			}
		}

		if pickedRoom == nil {
			return bookingFailedMsg{err: fmt.Errorf("no room suiting user's selection, aborting")}
		}

		if user.Credits < pickedRoom.Price {
			return bookingFailedMsg{err: fmt.Errorf("not enough credits to perform the booking, aborting")}
		}

		dt := b.getStartTime(b.browsePayload.StartDate, b.browsePayload.StartHour)

		payload := api.CosoftBookingPayload{
			CosoftAvailabilityPayload: api.CosoftAvailabilityPayload{
				DateTime: dt,
				NbPeople: b.browsePayload.NbPeople,
				Duration: b.browsePayload.Duration,
			},
			UserCredits: user.Credits,
			Room:        *pickedRoom,
		}

		apiClient := api.NewApi()

		err = apiClient.BookRoom(user.WAuth, user.WAuthRefresh, payload)

		if err != nil {
			return bookingFailedMsg{err: err}
		}

		return bookingCompleteMsg{room: *pickedRoom}
	}
}

func (b *BrowseModel) generateTable() string {
	dt := b.getStartTime(b.browsePayload.StartDate, b.browsePayload.StartHour)
	endTime := dt.Add(time.Duration(b.browsePayload.Duration) * time.Minute)
	dateFormat := "02/01/2006 15:04"

	paidPrice := b.bookedRoom.Price * (float64(b.browsePayload.Duration) / 60)

	headers := []string{"ROOM", "DURATION", "COST"}

	rows := [][]string{
		{
			b.bookedRoom.Name,
			fmt.Sprintf("%s → %s", dt.Format(dateFormat), endTime.Format(dateFormat)),
			fmt.Sprintf("%.2f credits", paidPrice),
		},
	}

	return common.CreateTable(headers, rows)
}

func (b *BrowseModel) getStartTime(startDate, startHour string) time.Time {
	// Already validated at this point.
	tDate, _ := time.Parse(time.DateOnly, startDate)
	tHour, _ := time.Parse("15:04", startHour)

	return time.Date(
		tDate.Year(),
		tDate.Month(),
		tDate.Day(),
		tHour.Hour(),
		tHour.Minute(),
		0,
		0,
		tDate.Location(),
	)
}

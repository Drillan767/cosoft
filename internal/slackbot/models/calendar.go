package models

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/storage"
	"cosoft-cli/shared/models"
	"fmt"
	"strings"
	"sync"
	"time"
)

type CalendarState struct {
	CurrentDate time.Time
	Calendar    string
	Error       *string
}

func newCalendarState() *CalendarState {
	return &CalendarState{
		CurrentDate: time.Now(),
	}
}

func (s *CalendarState) Type() string { return calendarStateType }

func (s *CalendarState) Update(store *storage.Store, params UpdateParams) (State, error) {
	switch params.ActionID {
	case "cancel":
		return NewLandingState(store, params.UserID)

	case "next-day":
		s.CurrentDate = s.CurrentDate.Add(24 * time.Hour)
		return s.load(store, params)

	case "prev-day":
		today := time.Now().Truncate(24 * time.Hour)
		if s.CurrentDate.Truncate(24 * time.Hour).Equal(today) {
			// already on today, can't go back
			return s, nil
		}
		s.CurrentDate = s.CurrentDate.Add(-24 * time.Hour)
		return s.load(store, params)

	default:
		return s, nil
	}
}

func (s *CalendarState) Next() bool { return false }

func (s *CalendarState) load(store *storage.Store, params UpdateParams) (State, error) {
	user, err := store.GetUserData(&params.UserID)
	if err != nil {
		// TODO: redirect the user to the login page and display an error?
		return s, fmt.Errorf("get user data: %v", err)
	}

	// Get user's future reservations
	apiClient := api.NewApi()
	bookings, err := apiClient.GetFutureBookings(user.WAuth, user.WAuthRefresh)
	if err != nil {
		s.Error = ptr(":red_circle: Impossible de charger les réservations")
		return s, fmt.Errorf("get future bookings: %v", err)
	}

	// Ensure we have all rooms available.
	rooms, err := listAndStoreRooms(store, user)
	if err != nil {
		s.Error = ptr(":red_circle: Impossible de récupérer les salles de réunion")
		return s, fmt.Errorf("list and store rooms: %v", err)
	}

	planning, err := getRoomsPlanning(
		user,
		rooms,
		s.CurrentDate,
		bookings.Data,
	)
	if err != nil {
		s.Error = ptr(":red_circle: Impossible de charger le calendrier")
	}

	s.Calendar = planning
	return s, nil
}

func listAndStoreRooms(store *storage.Store, user *storage.User) ([]storage.Room, error) {
	rooms, err := store.GetRooms()
	if err != nil {
		return nil, fmt.Errorf("get rooms: %v", err)
	}

	// Rooms are stored, return early
	if len(rooms) > 0 {
		return rooms, nil
	}

	// Fetch the rooms
	apiClient := api.NewApi()
	apiRooms, err := apiClient.GetAllRooms(user.WAuth, user.WAuthRefresh)
	if err != nil {
		return nil, fmt.Errorf("get all rooms: %v", err)
	}

	// Store them
	err = store.CreateRooms(apiRooms)
	if err != nil {
		return nil, fmt.Errorf("create rooms: %v", err)
	}

	rooms = make([]storage.Room, len(apiRooms))
	for i, room := range apiRooms {
		rooms[i] = storage.Room{
			Id:       room.Id,
			Name:     room.Name,
			MaxUsers: room.NbUsers,
			Price:    room.Price,
		}
	}

	// And return them
	return rooms, nil
}

func getRoomsPlanning(
	user *storage.User,
	rooms []storage.Room,
	date time.Time,
	userBookings []api.Reservation,
) (string, error) {
	apiClient := api.NewApi()

	// TODO: use error group to sync goroutines and stop on first error.
	var wg sync.WaitGroup
	results := make([]models.RoomUsage, len(rooms))
	for i, r := range rooms {
		wg.Go(func() {
			defer wg.Done()
			response, err := apiClient.GetRoomBusyTime(
				user.WAuth,
				user.WAuthRefresh,
				r.Id,
				date,
			)

			result := models.RoomUsage{
				Id:   r.Id,
				Name: r.Name,
			}

			if err == nil && response != nil {
				result.UsedSlots = *response
			} else {
				if err != nil {
					fmt.Printf("Error for %s: %s", r.Name, err.Error())
				} else {
					fmt.Printf("Nil response for %s", r.Name)
				}
			}

			results[i] = result
		})
	}
	wg.Wait()

	rows := common.BuildCalendar(0, 16, results, userBookings)
	return strings.Join(rows, "\n"), nil
}

package services

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/storage"
	"cosoft-cli/shared/models"
	"fmt"
	"sync"
	"time"
)

func (s *SlackService) getRoomAvailabilities(
	user storage.User,
	nbPeople, duration int,
	dateTime time.Time,
) ([]models.Room, error) {

	apiClient := api.NewApi()

	payload := api.CosoftAvailabilityPayload{
		DateTime: dateTime,
		NbPeople: nbPeople,
		Duration: duration,
	}

	rooms, err := apiClient.GetAvailableRooms(user.WAuth, user.WAuthRefresh, payload)

	if err != nil {
		return nil, err
	}

	if len(rooms) == 0 {
		return nil, fmt.Errorf(":red_circle: Aucune salle disponible")
	}

	return rooms, nil
}

func (s *SlackService) bookRoom(
	user storage.User,
	nbPeople, duration int,
	pickedRoom models.Room,
	dateTime time.Time,
) error {

	payload := api.CosoftBookingPayload{
		CosoftAvailabilityPayload: api.CosoftAvailabilityPayload{
			NbPeople: nbPeople,
			Duration: duration,
			DateTime: dateTime,
		},
		UserCredits: user.Credits,
		Room:        pickedRoom,
	}

	apiClient := api.NewApi()

	err := apiClient.BookRoom(user.WAuth, user.WAuthRefresh, payload)

	if err != nil {
		return err
	}

	return nil
}

func (s *SlackService) fetchReservations(user storage.User) ([]api.Reservation, error) {
	apiClient := api.NewApi()
	bookings, err := apiClient.GetFutureBookings(user.WAuth, user.WAuthRefresh)

	if err != nil {
		return nil, err
	}

	return bookings.Data, nil
}

func (s *SlackService) cancelReservation(
	user storage.User,
	reservationId string,
) error {
	apiClient := api.NewApi()

	return apiClient.CancelBooking(user.WAuth, user.WAuthRefresh, reservationId)
}

func (s *SlackService) getAllRooms(user storage.User) ([]storage.Room, error) {
	rooms, err := s.store.GetRooms()
	if err != nil {
		return nil, err
	}

	// Rooms are stored, return early
	if len(rooms) > 0 {
		return rooms, err
	}

	// Fetch the rooms
	apiClient := api.NewApi()
	apiRooms, err := apiClient.GetAllRooms(user.WAuth, user.WAuthRefresh)

	if err != nil {
		return nil, err
	}

	// Store them
	err = s.store.CreateRooms(apiRooms)

	if err != nil {
		return nil, err
	}

	rooms = make([]storage.Room, 0, len(apiRooms))

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

func (s *SlackService) getRoomsPlanning(
	user *storage.User,
	rooms []storage.Room,
	date time.Time,
	userBookings []api.Reservation,
) (string, error) {
	apiClient := api.NewApi()
	results := make([]models.RoomUsage, len(rooms))
	var wg sync.WaitGroup

	for i, r := range rooms {
		wg.Add(1)
		go func(i int, r storage.Room) {
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
		}(i, r)
	}

	wg.Wait()

	rows := common.BuildCalendar(0, 16, results, userBookings)

	var calendar string

	for _, row := range rows {
		calendar = fmt.Sprintf("%s\n%s", calendar, row)
	}

	return calendar, nil
}

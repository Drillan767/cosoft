package services

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/common"
	"cosoft-cli/internal/storage"
	"cosoft-cli/shared/models"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (s *Service) UpdateCredits() (*float64, error) {
	return s.store.UpdateCredits(nil)
}

func (s *Service) EnsureRoomsStored() error {
	rooms, err := s.store.GetRooms()

	if err != nil {
		return err
	}

	if len(rooms) > 0 {
		return nil
	}

	authData, err := s.store.GetUserData(nil)

	if err != nil {
		return err
	}

	apiClient := api.NewApi()
	apiRooms, err := apiClient.GetAllRooms(authData.WAuth, authData.WAuthRefresh)

	if err != nil {
		return err
	}

	return s.store.CreateRooms(apiRooms)
}

func (s *Service) GetRoomAvailabilities(
	date time.Time,
	userBookings []api.Reservation,
) ([]string, error) {
	rooms, err := s.store.GetRooms()

	if err != nil {
		return nil, err
	}

	authData, err := s.store.GetUserData(nil)

	if err != nil {
		return nil, err
	}

	apiClient := api.NewApi()

	results := make([]models.RoomUsage, len(rooms))
	var wg sync.WaitGroup

	for i, room := range rooms {
		wg.Add(1)
		go func(i int, room storage.Room) {
			defer wg.Done()
			response, err := apiClient.GetRoomBusyTime(
				authData.WAuth,
				authData.WAuthRefresh,
				room.Id,
				date,
			)

			result := models.RoomUsage{
				Id:   room.Id,
				Name: room.Name,
			}

			if err == nil && response != nil {
				result.UsedSlots = *response
			} else {
				if err != nil {
					debug(fmt.Sprintf("Error for %s: %s", room.Name, err.Error()))
				} else {
					debug(fmt.Sprintf("Nil response for %s", room.Name))
				}
			}

			results[i] = result
		}(i, room)
	}

	wg.Wait()

	rows := common.BuildCalendar(0, 16, results, userBookings)

	return rows, nil
}

func (s *Service) NonInteractiveBooking(
	capacity, duration int,
	name string,
	dt time.Time,
) (string, error) {
	user, err := s.store.GetUserData(nil)

	if err != nil {
		return "", err
	}

	clientApi := api.NewApi()

	// Ensure user is authenticated
	fmt.Println("checking user authentication status...")
	err = clientApi.GetAuth(user.WAuth, user.WAuthRefresh)

	if err != nil {
		return "", fmt.Errorf("user not authenticated: %v", err)
	}

	var room *models.Room

	payload := api.CosoftAvailabilityPayload{
		DateTime: dt,
		Duration: duration,
		NbPeople: capacity,
	}

	fmt.Println("retrieving available rooms with requested filters...")
	availabilities, err := clientApi.GetAvailableRooms(user.WAuth, user.WAuthRefresh, payload)

	if err != nil {
		return "", err
	}

	if len(availabilities) == 0 {
		return "", errors.New("no available rooms")
	}

	// If room name was provided, check if is among the API's response.
	if name != "" {
		var found *models.Room
		for _, avail := range availabilities {
			if avail.Name == name {
				found = &avail
				break
			}
		}

		if found == nil {
			return "", fmt.Errorf("room %s not available for the selected filter", name)
		}

		room = found
	}

	// Set room id as either the asked room's id, or the 1st available room id
	targetRoom := availabilities[0]

	if room != nil {
		targetRoom = *room
	}

	if targetRoom.Price > user.Credits {
		return "", errors.New("not enough credits")
	}

	bookingPayload := api.CosoftBookingPayload{
		CosoftAvailabilityPayload: api.CosoftAvailabilityPayload{
			DateTime: dt,
			Duration: duration,
			NbPeople: capacity,
		},
		UserCredits: user.Credits,
		Room:        targetRoom,
	}

	fmt.Println("booking requested room...")
	err = clientApi.BookRoom(user.WAuth, user.WAuthRefresh, bookingPayload)

	if err != nil {
		return "", err
	}

	endTime := dt.Add(time.Duration(duration) * time.Minute)
	success := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render(`✓ Booking complete!`)
	dateFormat := "02/01/2006 15:04"

	headers := []string{"ROOM", "DURATION", "COST"}

	rows := [][]string{
		{
			targetRoom.Name,
			fmt.Sprintf("%s → %s", dt.Format(dateFormat), endTime.Format(dateFormat)),
			fmt.Sprintf("%.2f credits", targetRoom.Price),
		},
	}

	fmt.Println(success)

	return common.CreateTable(headers, rows), nil
}

func debug(text string) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s: %s \n\n", time.Now().Format(time.RFC3339), text))
}

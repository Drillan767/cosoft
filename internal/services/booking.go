package services

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
	"cosoft-cli/shared/models"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func (s *Service) UpdateCredits() (*float64, error) {
	return s.store.UpdateCredits()
}

func (s *Service) EnsureRoomsStored() error {
	rooms, err := s.store.GetRooms()

	if err != nil {
		return err
	}

	if len(rooms) > 0 {
		return nil
	}

	authData, err := s.store.GetUserData()

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

func (s *Service) GetRoomAvailabilities(date time.Time) ([]string, error) {
	rooms, err := s.store.GetRooms()

	if err != nil {
		return nil, err
	}

	authData, err := s.store.GetUserData()

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
				debug(fmt.Sprintf("Got response for %s", room.Name))
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

	maxLabelLength := 0

	for _, room := range results {
		if len(room.Name)+1 > maxLabelLength {
			maxLabelLength = len(room.Name) + 1
		}
	}

	rows := make([]string, len(results)+1)

	rows[0] = s.createTableHeader(maxLabelLength)

	for i, room := range results {
		rows[i+1] = s.createTableRow(room, maxLabelLength)
	}

	return rows, nil
}

func (s *Service) createTableHeader(labelLength int) string {
	spacing := 2
	displayedHours := 16
	result := ""

	for i := 0; i < displayedHours; i++ {
		if i+8 < 10 {
			result += "0"
		}
		result += fmt.Sprintf("%dh%s", i+8, strings.Repeat(" ", spacing))
	}

	return strings.Repeat(" ", labelLength-1) + result
}

func (s *Service) createTableRow(row models.RoomUsage, labelLength int) string {
	spacing := labelLength - len(row.Name)
	return row.Name + strings.Repeat(" ", spacing) + "│"
}

func debug(text string) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s: %s \n\n", time.Now().Format(time.RFC3339), text))
}

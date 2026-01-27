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
	displayedHours := 16

	for _, room := range results {
		if len(room.Name)+1 > maxLabelLength {
			maxLabelLength = len(room.Name) + 1
		}
	}

	rows := make([]string, len(results)+1)

	rows[0] = s.createTableHeader(maxLabelLength, displayedHours)

	for i, room := range results {
		rows[i+1] = s.createTableRow(room, maxLabelLength, displayedHours)
	}

	return rows, nil
}

func (s *Service) createTableHeader(labelLength, displayedHours int) string {
	spacing := 2
	result := ""

	for i := 0; i < displayedHours; i++ {
		if i+8 < 10 {
			result += "0"
		}
		result += fmt.Sprintf("%dh%s", i+8, strings.Repeat(" ", spacing))
	}

	return strings.Repeat(" ", labelLength-1) + result
}

func (s *Service) createTableRow(row models.RoomUsage, labelLength, displayedHours int) string {
	type parsedSlot struct {
		Start time.Time
		End   time.Time
	}

	var slots []parsedSlot
	spacing := labelLength - len(row.Name)
	columns := ""

	for _, slot := range row.UsedSlots {
		debug(fmt.Sprintf("%+v\n", slot))
		start, _ := time.Parse("2006-01-02T15:04:05", slot.Start)
		end, _ := time.Parse("2006-01-02T15:04:05", slot.End)

		slots = append(slots, parsedSlot{
			Start: start,
			End:   end,
		})
	}

	// now := time.Now()
	baseDate := slots[0].Start.Truncate(24 * time.Hour)
	startTime := baseDate.Add(8 * time.Hour)
	endTime := baseDate.Add(23 * time.Hour)
	counter := 0

	current := startTime

	for !current.After(endTime) {
		occupied := false
		slotEnd := current.Add(15 * time.Minute)

		for _, slot := range slots {
			if current.Before(slot.End) && slotEnd.After(slot.Start) {
				occupied = true
				break
			}
		}

		symbol := " "

		if occupied {

			symbol = "█"
		}

		if counter%4 == 3 {
			symbol += "│"
		}

		columns += symbol
		counter++

		current = current.Add(15 * time.Minute)
	}

	return row.Name + strings.Repeat(" ", spacing) + "│" + columns
}

func debug(text string) {
	file, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s: %s \n\n", time.Now().Format(time.RFC3339), text))
}

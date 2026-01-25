package services

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
	"cosoft-cli/shared/models"
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

func (s *Service) GetRoomAvailabilities(date time.Time) ([]models.RoomUsage, error) {
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
			}

			results[i] = result
		}(i, room)
	}

	wg.Wait()

	return results, nil
}

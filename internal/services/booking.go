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

func (s *Service) GetRoomAvailabilities(date time.Time) (*models.RoomUsage, error) {
	rooms, err := s.store.GetRooms()

	if err != nil {
		return nil, err
	}

	authData, err := s.store.GetUserData()

	if err != nil {
		return nil, err
	}

	apiClient := api.NewApi()

	results := make([]storage.Room, len(rooms))
	var wg sync.WaitGroup

	/*
			 for i, room := range rooms {
			  wg.Add(1)
			  go func(idx int, r Room) {
				  defer wg.Done()
				  results[idx] = fetchBusyTimes(r) // each goroutine writes to its own index
			  }(i, room)
		  }
		  wg.Wait()
	*/

	for i, room := range rooms {
		wg.Add(1)
		go func(i int, room storage.Room) {
			defer wg.Done()
			results[i] = apiClient.GetAllAvailabilities(
				authData.WAuth,
				authData.WAuthRefresh,
				room.Id,
				date,
			)
		}(i, room)
	}

	wg.Wait()

	return nil, nil
}

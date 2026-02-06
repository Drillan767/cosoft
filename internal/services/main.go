package services

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
	"fmt"
	"os"
	"path/filepath"
)

type Service struct {
	store *storage.Store
}

func NewService() (*Service, error) {
	configDir, _ := os.UserConfigDir()
	store, err := storage.NewStore(fmt.Sprintf("%s/cosoft/data.db", configDir))

	if err != nil {
		return nil, err
	}

	return &Service{
		store: store,
	}, nil
}

func (s *Service) ClearData() error {
	// Disconnect user from api
	user, err := s.GetAuthData()

	if err != nil {
		return err
	}

	clientApi := api.NewApi()
	err = clientApi.Logout(user.WAuth, user.WAuthRefresh)

	if err != nil {
		return err
	}

	s.store.Close()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	cosoftDir := filepath.Join(configDir, "cosoft")

	return os.RemoveAll(cosoftDir)
}

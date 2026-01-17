package services

import (
	"cosoft-cli/internal/storage"
	"fmt"
	"os"
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

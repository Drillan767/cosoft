package services

import (
	"cosoft-cli/internal/storage"
)

type SlackService struct {
	store *storage.Store
}

func NewSlackService() (*SlackService, error) {
	store, err := storage.NewStore("./slack/database.db")

	if err != nil {
		return nil, err
	}

	return &SlackService{store}, nil
}

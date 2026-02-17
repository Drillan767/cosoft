package services

import (
	"cosoft-cli/internal/storage"
)

type SlackService struct {
	store *storage.Store
}

func NewSlackService(store *storage.Store) *SlackService {
	return &SlackService{store}
}

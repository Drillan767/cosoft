package auth

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
	"fmt"
	"os"
	"time"
)

type AuthService struct {
	store *storage.Store
}

func NewAuthService() (*AuthService, error) {
	configDir, _ := os.UserConfigDir()
	store, err := storage.NewStore(fmt.Sprintf("%s/cosoft/data.db", configDir))

	if err != nil {
		return nil, err
	}

	return &AuthService{
		store: store,
	}, nil
}

func (a *AuthService) IsAuthenticated() bool {

	hasActiveToken, err := a.store.HasActiveToken()

	if err != nil {
		return false
	}

	return hasActiveToken
}

func (a *AuthService) SaveAuthData(user *api.UserResponse, expiresAt time.Time) error {
	return a.store.CreateUser(user, expiresAt)
}

func (a *AuthService) Logout() error {
	return a.store.LogoutUser()
}

func (a *AuthService) GetAuthData() (*storage.User, error) {
	user, err := a.store.GetUserData()

	if err != nil {
		return nil, err
	}

	return user, nil
}

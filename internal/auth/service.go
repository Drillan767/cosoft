package auth

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
	"fmt"
	"os"
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

func (a *AuthService) SaveAuthData(user *api.UserResponse) error {
	return a.store.CreateUser(user, user.JwtToken, user.RefreshToken)
}

func (a *AuthService) Logout() error {
	return a.store.LogoutUser()
}

func (a *AuthService) GetAuthData() (*storage.User, error) {
	return a.store.GetUserData()
}

func (a *AuthService) UpdateCredits() (*float64, error) {
	return a.store.UpdateCredits()
}

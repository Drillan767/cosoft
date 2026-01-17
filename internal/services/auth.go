package services

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
)

func (s *Service) IsAuthenticated() bool {

	hasActiveToken, err := s.store.HasActiveToken()

	if err != nil {
		return false
	}

	return hasActiveToken
}

func (s *Service) SaveAuthData(user *api.UserResponse) error {
	return s.store.CreateUser(user, user.JwtToken, user.RefreshToken)
}

func (s *Service) Logout() error {
	return s.store.LogoutUser()
}

func (s *Service) GetAuthData() (*storage.User, error) {
	return s.store.GetUserData()
}

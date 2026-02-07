package services

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/storage"
)

func (s *Service) IsAuthenticated() bool {

	cookies, err := s.store.HasActiveToken(nil)

	if err != nil || cookies == nil {
		return false
	}

	apiClient := api.NewApi()
	err = apiClient.GetAuth(cookies.WAuth, cookies.WAuthRefresh)

	return err == nil
}

func (s *Service) SaveAuthData(user *api.UserResponse) error {
	return s.store.SetUser(user, user.JwtToken, user.RefreshToken, nil)
}

func (s *Service) Logout() error {
	return s.store.LogoutUser(nil)
}

func (s *Service) GetAuthData() (*storage.User, error) {
	return s.store.GetUserData(nil)
}

package auth

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/settings"
	"encoding/json"
	"os"
	"time"
)

type AuthService struct {
	configPath string
}

func NewAuthService() *AuthService {
	configDir, _ := os.UserConfigDir()
	return &AuthService{
		configPath: configDir + "/cosoft",
	}
}

func (a *AuthService) IsAuthenticated() bool {
	tokenPath := a.configPath + "/auth_data.json"

	data, err := os.ReadFile(tokenPath)

	if err != nil {
		return false
	}

	var ad settings.AuthData

	if err := json.Unmarshal(data, &ad); err != nil {
		return false
	}

	// Check token exists
	if ad.Token == "" {
		return false
	}

	// Check if token is not expired
	if time.Now().After(ad.ExpirationDate) {
		return false
	}

	return true
}

func (a *AuthService) SaveAuthData(user *api.UserResponse, expiresAt time.Time) error {
	tokenPath := a.configPath + "/auth_data.json"

	ad := settings.AuthData{
		Token:          user.JwtToken,
		ExpirationDate: expiresAt,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Credits:        user.Credits,
	}

	data, err := json.Marshal(ad)

	if err != nil {
		return err
	}

	return os.WriteFile(tokenPath, data, 0600)
}

func (a *AuthService) GetToken() (string, error) {
	// TODO: make the fetch
	return "", nil
}

func (a *AuthService) Logout() error {
	tokenPath := a.configPath + "/auth_data.json"

	ad := settings.AuthData{
		Token:          "",
		ExpirationDate: time.Now(),
	}

	data, _ := json.Marshal(ad)

	return os.WriteFile(tokenPath, data, 0600)
}

func (a *AuthService) GetAuthData() (*settings.AuthData, error) {
	tokenPath := a.configPath + "/auth_data.json"

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}

	var ad settings.AuthData
	if err := json.Unmarshal(data, &ad); err != nil {
		return nil, err
	}

	return &ad, nil
}

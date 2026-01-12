package auth

import (
	"cosoft-cli/internal/settings"
	"cosoft-cli/internal/ui"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type AuthService struct {
	configPath string
}

func NewAuthService() *AuthService {
	return &AuthService{
		configPath: os.Getenv("HOME") + "/.cosoft",
	}
}

func (a *AuthService) IsAuthenticated() bool {
	tokenPath := a.configPath + "/jwt_token.json"

	data, err := os.ReadFile(tokenPath)

	if err != nil {
		return false
	}

	var rt settings.JwtToken

	if err := json.Unmarshal(data, &rt); err != nil {
		return false
	}

	// Check token exists
	if rt.Token == "" {
		return false
	}

	// Check if token is not expired
	if time.Now().After(rt.ExpirationDate) {
		return false
	}

	return true
}

func (a *AuthService) RequiresAuth() error {
	if a.IsAuthenticated() {
		return nil
	}

	// Not authenticated, show form
	ui := ui.NewUI()
	loginModel, err := ui.LoginForm()

	if err != nil {
		return err
	}

	token := loginModel.GetToken()

	// Check if token is actually present (login succeeded)
	if token == "" {
		return fmt.Errorf("authentication cancelled or failed")
	}

	// Todo: replace with actual expiration date from actual token
	expirationDate := time.Now().Add(7 * 24 * time.Hour) // 1 week

	return a.SaveToken(token, expirationDate)
}

func (a *AuthService) SaveToken(token string, expiresAt time.Time) error {
	tokenPath := a.configPath + "/jwt_token.json"

	rt := settings.JwtToken{
		Token:          token,
		ExpirationDate: expiresAt,
	}

	data, err := json.Marshal(rt)

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
	tokenPath := a.configPath + "/jwt_token.json"

	rt := settings.JwtToken{
		Token:          "",
		ExpirationDate: time.Now(),
	}

	data, _ := json.Marshal(rt)

	return os.WriteFile(tokenPath, data, 0600)
}

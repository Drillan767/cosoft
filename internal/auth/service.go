package auth

import (
	"cosoft-cli/internal/api"
	"cosoft-cli/internal/settings"
	"cosoft-cli/internal/ui"
	"encoding/json"
	"fmt"
	"os"
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
	tokenPath := a.configPath + "/refresh_token.json"

	data, err := os.ReadFile(tokenPath)

	if err != nil {
		return false
	}

	var rt settings.RefreshToken

	if err := json.Unmarshal(data, &err); err != nil {
		return false
	}

	return rt.RefreshToken != ""
}

func (a *AuthService) RequiresAuth() error {
	if a.IsAuthenticated() {
		return nil
	}

	// Not authenticated, show form
	ui := ui.NewUI()
	creds, err := ui.LoginFormWithLayout()
	if err != nil {
		return err
	}

	api := api.NewApi()
	err = api.Login(creds) // TODO: return actual token

	if err != nil {
		return fmt.Errorf("Login failed: %w", err)
	}

	return a.SaveToken("token")

}

func (a *AuthService) SaveToken(token string) error {
	tokenPath := a.configPath + "/refresh_token.json"

	rt := settings.RefreshToken{RefreshToken: token}

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
	tokenPath := a.configPath + "/refresh_token.json"

	rt := settings.RefreshToken{RefreshToken: ""}

	data, _ := json.Marshal(rt)

	return os.WriteFile(tokenPath, data, 0600)
}

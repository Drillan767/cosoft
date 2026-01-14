package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type UserConfig struct {
	FavoriteRooms    []string `json:"favoriteRooms"`
	PreferedDuration int      `json:"preferedDuration"`
}

type AuthData struct {
	Token          string    `json:"token"`
	ExpirationDate time.Time `json:"expires"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	Credits        float64   `json:"credits"`
}

func LoadConfiguration() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	// Use "cosoft" as the application directory name
	path := fmt.Sprintf("%s/cosoft", configDir)

	err = createConfiguration(path)

	if err != nil {
		return err
	}

	err = ensureUserSettingsExist(path)

	if err != nil {
		return err
	}

	err = ensureAuthDataFileExists(path)

	if err != nil {
		return err
	}

	return nil
}

func createConfiguration(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Configuration directory does not exist, creating...")

		err = os.MkdirAll(path, 0700)

		if err != nil {
			return err
		}
	}

	return nil
}

func ensureUserSettingsExist(path string) error {
	settingsPath := fmt.Sprintf("%s/user_settings.json", path)

	if _, err := os.Stat(settingsPath); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("User settings do not exist, creating...")
		userConfig := UserConfig{
			FavoriteRooms:    []string{},
			PreferedDuration: 60,
		}

		jsonValue, err := json.Marshal(userConfig)

		if err != nil {
			return err
		}

		err = os.WriteFile(settingsPath, jsonValue, 0640)

		if err != nil {
			return err
		}

		return nil

	} else {
		return err
	}
}

func ensureAuthDataFileExists(path string) error {
	authDataPath := fmt.Sprintf("%s/auth_data.json", path)

	if _, err := os.Stat(authDataPath); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Auth data does not exist, creating...")
		ad := AuthData{
			Token:          "",
			ExpirationDate: time.Now(),
		}

		adValue, err := json.Marshal(ad)

		if err != nil {
			return err
		}

		err = os.WriteFile(authDataPath, adValue, 0640)

		if err != nil {
			return err
		}

		return nil

	} else {
		return err
	}
}

package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type UserConfig struct {
	FavoriteRooms    []string `json:"favoriteRooms"`
	PreferedDuration int      `json:"preferedDuration"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}

func LoadConfiguration() error {
	path := fmt.Sprintf("%s/.cosoft", os.Getenv("HOME"))

	err := createConfiguration(path)

	if err != nil {
		return err
	}

	err = ensureUserSettingsExist(path)

	if err != nil {
		return err
	}

	err = ensureRefreshTokenFileExists(path)

	if err != nil {
		return err
	}

	return nil
}

func createConfiguration(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Configuration directory does not exist, creating...")

		err = os.Mkdir(path, 0700)

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

func ensureRefreshTokenFileExists(path string) error {
	refreshTokenPath := fmt.Sprintf("%s/refresh_token.json", path)

	if _, err := os.Stat(refreshTokenPath); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Refresh token does not exist, creating...")
		rt := RefreshToken{
			RefreshToken: "",
		}

		rtValue, err := json.Marshal(rt)

		if err != nil {
			return err
		}

		err = os.WriteFile(refreshTokenPath, rtValue, 0640)

		if err != nil {
			return err
		}

		return nil

	} else {
		return err
	}
}

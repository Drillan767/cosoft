package settings

import (
	"cosoft-cli/internal/storage"
	"fmt"
	"os"
)

type UserConfig struct {
	FavoriteRooms    []string `json:"favoriteRooms"`
	PreferedDuration int      `json:"preferedDuration"`
}

func EnsureDatabaseExists() error {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return err
	}

	cosoftDir := fmt.Sprintf("%s/cosoft", configDir)
	path := fmt.Sprintf("%s/data.db", cosoftDir)

	// Ensure the cosoft directory exists
	if err := os.MkdirAll(cosoftDir, 0755); err != nil {
		return err
	}

	store, err := storage.NewStore(path)

	if err != nil {
		return err
	}

	err = store.SetupDatabase(path)

	if err != nil {
		return err
	}

	return nil
}

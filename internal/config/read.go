package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}

	var configStruct Config
	if unmarshalErr := json.Unmarshal(file, &configStruct); unmarshalErr != nil {
		return Config{}, unmarshalErr
	}
	return configStruct, nil
}

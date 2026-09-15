package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	fullPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

func (c *Config) SetUser(name string) error {
	c.UserName = name
	err := write(*c)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_URL       string `json:"db_url"`
	CURRENT_USER string `json:"current_user_name"`
}

func Read() (Config, error) {
	var config Config
	home, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(home)
	if err != nil {
		return Config{}, err
	}
	if err = json.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}
	return config, nil

}

func (c *Config) SetUser(userName string) error {
	c.CURRENT_USER = userName
	return c.write()
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + "/" + configFileName, nil
}

func (c *Config) write() error {
	home, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(home, file, 0644)

}

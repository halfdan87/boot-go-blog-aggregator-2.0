package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DbUrl         string `json:"db_url"`
	CurrentUserId string `json:"current_user_id"`
}

func Read() (*Config, error) {
	jsonFile, err := os.Open(getConfigPath())
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (cfg *Config) SetUser(uname string) error {
	cfg.CurrentUserId = uname

	jsonFile, err := os.OpenFile(getConfigPath(), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := json.Marshal(cfg)
	_, err = jsonFile.Write(byteValue)
	if err != nil {
		return err
	}
	return nil
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return home + string(os.PathSeparator) + ".gatorconfig.json"
}

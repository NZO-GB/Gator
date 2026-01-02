package internal

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"
const path := os.UserHomeDir() + configFileName

type Config struct{
	DbUrl: string `json:"db_url"`
	CurrentUserName: string `json:"CurrentUserName"`
}

func Read() (Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func write(cfg Config) error {

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data)
}

func (c Config) SetUser(username string) {
	cfg := Read()
	cfg.CurrentUserName = username
	write(cfg)
}
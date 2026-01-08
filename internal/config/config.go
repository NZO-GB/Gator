package internal

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl			string `json:"db_url"`
	CurrentUserName	string `json:"CurrentUserName"`
	Protocol		string `json:"protocol"`
}

func Read() (Config, error) {

	homedir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	path := homedir + "/" + configFileName

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func write(cfg Config) error {

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := homedir + "/" + configFileName

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c Config) SetUser(username string) error {
	cfg, err := Read()
	if err != nil {
		return err
	}
	cfg.CurrentUserName = username
	err = write(cfg)
	if err != nil {
		return err
	}
	return nil
}

func (c Config) AddProtocol(connectionString string) error {
	cfg, err := Read()
	if err != nil {
		return err
	}
	cfg.Protocol = connectionString
	err = write(cfg)
	if err != nil {
		return err
	}
	return nil
}
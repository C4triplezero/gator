package config

import (
	"encoding/json"
	"os"
)

func (cfg Config) SetUser(username string) error {
	cfg.CurrentUserName = username

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, os.FileMode(0666))
	if err != nil {
		return err
	}

	return nil
}

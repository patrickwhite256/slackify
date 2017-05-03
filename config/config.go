package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	SlackClientId       string `json:"slack_client_id"`
	SlackClientSecret   string `json:"slack_client_secret"`
	SpotifyClientId     string `json:"spotify_client_id"`
	SpotifyClientSecret string `json:"spotify_client_secret"`

	DBAddress  string `json:"db_address"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBDatabase string `json:"db_database"`
}

func FromFile(filename string) (*Config, error) {
	conf := &Config{}
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileData, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

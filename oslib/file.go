package oslib

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Users  []string `json:"users"`
	Orgs   []string `json:"orgs"`
	Labels []string `json:"labels"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

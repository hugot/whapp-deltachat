package core

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	DataFolder  string `yaml:"data_folder"`
	UserAddress string `yaml:"user"`
	ShowFromMe  bool   `yaml:"show_from_me"`
}

type Config struct {
	Deltachat map[string]string `yaml:"deltachat"`
	App       AppConfig         `yaml:"app"`
}

func ConfigFromFile(filePath string) (*Config, error) {
	fileContents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	config := &Config{}

	return config, yaml.Unmarshal(fileContents, config)
}

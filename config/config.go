package config

import (
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		Token  string  `yaml:"token"`
		Admins []int64 `yaml:"admins"`
	} `yaml:"bot"`

	Itemshop struct {
		Channel           int64  `yaml:"channel"`
		PrimaryTitle      string `yaml:"primary_title"`
		PrimaryPrice      int    `yaml:"primary_price"`
		PrimaryCurrency   string `yaml:"primary_currency"`
		SecondaryTitle    string `yaml:"secondary_title"`
		SecondaryPrice    int    `yaml:"secondary_price"`
		SecondaryCurrency string `yaml:"secondary_currency"`
	} `yaml:"itemshop"`
}

func Load() (*Config, error) {
	file, err := os.Open("./config.yaml")
	if err != nil {
		return nil, err
	}

	conf := &Config{}

	if err = yaml.NewDecoder(file).Decode(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (data *Config) Save() error {
	configBytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile("./config.yaml", configBytes, fs.ModePerm)
}

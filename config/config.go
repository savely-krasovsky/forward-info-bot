package config

import (
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	LogLevel string `yaml:"LogLevel" validate:"required"`

	*Telegram `yaml:"Telegram" validate:"required"`
}

type Telegram struct {
	Token string `yaml:"Token" validate:"required"`
}

// Init new config with validation
func NewConfig(p string) (*Config, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

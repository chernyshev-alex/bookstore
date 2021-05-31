package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Application configuration

type Config struct {
	Elastic struct {
		URL string `yaml:"url" env:"ELK_URL" env-default:"http://127.0.0.1:9200"`
	} `yaml:"elastic"`
}

func LoadCondigFromFile(configPath string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

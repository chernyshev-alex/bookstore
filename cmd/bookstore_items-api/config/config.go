package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Application configuration
// TODO consider https://github.com/spf13/viper

type Config struct {
	Elastic struct {
		URL string `yaml:"url" env:"ELK_URL" env-default:"http://127.0.0.1:9200"`
	} `yaml:"elastic"`

	Server struct {
		Host string `yaml:"host" env:"SRV_HOST,HOST" env-default:"localhost"`
		Port string `yaml:"port" env:"SRV_PORT,PORT" env-default:"8081"`
	} `yaml:"server"`
}

func LoadCondigFromFile(configPath string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

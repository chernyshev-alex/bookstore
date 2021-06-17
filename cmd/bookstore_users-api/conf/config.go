package conf

import (
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

// Application configuration

type Config struct {
	Database struct {
		Host   string `yaml:"host" env:"DB_HOST" env-description:"db host"`
		Port   string `yaml:"port" env:"DB_PORT"`
		Schema string `yaml:"schema" env:"SCHEMA"`
		Uname  string `yaml:"uname" env:"USER_NAME"`
	} `yaml:"database"`

	OAuth struct {
		URL string `yaml:"url" env:"OAUT_URL" env-default:"http://127.0.0.1:8082"`
	} `yaml:"oauth"`

	Server struct {
		Host string `yaml:"host" env:"SRV_HOST,HOST" env-default:"localhost"`
		Port string `yaml:"port" env:"SRV_PORT,PORT" env-default:"8081"`
	} `yaml:"server"`
}

func LoadConfigFromEnv(envVar string) (*Config, error) {
	configFile, _ := filepath.Abs(os.Getenv(envVar))
	return LoadConfigFromFile(configFile)
}

func LoadConfigFromFile(configPath string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

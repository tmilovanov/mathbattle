package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TelegramToken            string `yaml:"telegram_token"`
	APIUrl                   string `yaml:"api_url"`
	DatabaseType             string `yaml:"db_type"`
	DatabaseConnectionString string `yaml:"db_connection_string"`
	ProblemsPath             string `yaml:"problems_path"`
	SolutionsPath            string `yaml:"solutions_path"`
}

func LoadConfig(configPath string) Config {
	result := Config{}
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config path, error: %v", err)
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&result)
	if err != nil {
		log.Fatalf("Failed to decode config, error: %v", err)
	}

	return result
}

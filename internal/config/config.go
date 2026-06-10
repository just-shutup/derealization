/// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"derealization/internal/bot"
)

type AppConfig struct {
	App struct {
		LogLevel string `yaml:"log_level"`
		LogFile  string `yaml:"log_file"`
		MaxBots  int    `yaml:"max_bots"`
	} `yaml:"app"`
	Bots []bot.BotConfig `yaml:"bots"` // Используем структуру из bot/config.go
}

// Load читает YAML файл и возвращает AppConfig.
func Load(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия конфига: %w", err)
	}
	defer func() { _ = file.Close() }()

	var cfg AppConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	return &cfg, nil
}
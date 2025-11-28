package config

import (
	"fmt"

	"github.com/sater-151/todo-list/internal/models"
	"github.com/spf13/viper"
)

const (
	configFile = "configuration.yaml"
)

var password string

func GetConfig() (models.Config, error) {
	var config models.Config

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("ошибка чтения конфигурации: %w", err)
	}

	// Преобразование конфигурации в структуру
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("ошибка преобразования конфигурации: %w", err)
	}

	password = config.HttpClient.Password

	return config, nil
}

func GetPass() (pass string) {
	pass = password
	return pass
}

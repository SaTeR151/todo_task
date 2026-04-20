package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sater-151/todo-list/pkg/utils"
)

type Config struct {
	Database struct {
		Host      string `env:"DB_HOST" env-required:"true"`
		Port      string `env:"DB_PORT" env-required:"true"`
		User      string `env:"DB_USER" env-required:"true"`
		Password  string `env:"DB_PASSWORD" env-required:"true"`
		DbName    string `env:"DB_NAME" env-required:"true"`
		CryproKey string `env:"DB_CRYPRO_KEY" env-required:"true"`
		Schema    string `env:"DB_SCHEMA" env-required:"true"`
	}

	SecretKey string `env:"SECRET_KEY" env-required:"true"`
}

func GetConfig() (cfg Config, err error) {
	defer utils.AddFuncLabel("[init-get-config]", err)

	if err = godotenv.Load(); err != nil {
		return
	}

	if err = cleanenv.ReadEnv(&cfg); err != nil {
		return
	}

	return
}

package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type StorageConfig struct {
	Username        string
	Password        string
	CryptoKey       string
	Host            string
	Port            string
	Database        string
	ApplicationName string
}

func NewPostgres(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg StorageConfig) (pool *pgxpool.Pool, err error) {
	if cfg.ApplicationName == "" {
		cfg.ApplicationName = "UnknownApp"
	}

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?application_name=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.ApplicationName,
	)

	logrus.Infof("Начинаю подключение к базе данных: postgresql://%s:***@%s:%s/%s?application_name=%s",
		cfg.Username, cfg.Host, cfg.Port, cfg.Database, cfg.ApplicationName,
	)

	err = DoWithAttempts(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Генерируем конфигурацию для подключения к postgres
		pgxCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return fmt.Errorf("не могу распознать конфигурацию для подключения к базе данных: %v", err)
		}

		// if logrusLogger != nil {
		// 	// Подключаем к конфигурации логгер logrus
		// 	pgxCfg.ConnConfig.Logger = logrusadapter.NewLogger(logrusLogger)
		// } else {
		// 	// Подключаем к конфигурации логгер, записывающий запросы и время их выполнения как трассировки
		// 	pgxCfg.ConnConfig.Logger = tracing.NewPgxLogger()
		// }

		// Подключаемся к postgres
		pool, err = pgxpool.NewWithConfig(ctx, pgxCfg)
		if err != nil {
			logrus.Error("Невозможно подключиться к базе данных... Выполняю следующую попытку...")
			return err
		}

		// Ping postgeres что бы удостовериться в доступности базы
		if err := pool.Ping(ctx); err != nil {
			logrus.Errorf("Не могу выполнить Ping к базе данных: %v", err)
			return err
		}

		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		return nil, fmt.Errorf("все попытки подключения к базе данных неудачны. Невозможно подключиться: %v", err)
	}

	logrus.Infof("Подключение к базе данных %s:%s успешно", cfg.Host, cfg.Port)

	return pool, nil
}

func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}

	return err
}

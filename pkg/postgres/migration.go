package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

type MigrationConfig struct {
	MigrationPath string
	MigrationName string
}

func BeginMigrations(cfg StorageConfig, mcfg MigrationConfig) error {

	// Если имя таблицы для хранения даных о версии миграции
	// не было задано, то оно будет установлено по умолчанию
	x_migrations_table := mcfg.MigrationName
	if x_migrations_table == "" {
		x_migrations_table = "schema_migrations"
	}

	if mcfg.MigrationPath == "" {
		return fmt.Errorf("не указан путь до файлов миграции. пример: file://./migrations")
	}

	destination := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, x_migrations_table,
	)

	logrus.Infof("Начинаю применять миграции к базе данных: postgres://%s:***@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Host, cfg.Port, cfg.Database,
	)

	m, err := migrate.New(mcfg.MigrationPath, destination)
	if err != nil {
		return fmt.Errorf("ошибка миграции: %v", err)
	}

	defer func() {
		logrus.Infof("[pkg/postgres] завершаю миграции к базе данных: postgres://%s:***@%s:%s/%s?sslmode=disable",
			cfg.Username, cfg.Host, cfg.Port, cfg.Database)
		if source_err, database_err := m.Close(); source_err != nil || database_err != nil {
			if err == nil { // не перетираем уже существующую ошибку
				err = fmt.Errorf("migration close source_err: %v  database_err: %v", source_err, database_err)
			}
		} else {
			logrus.Infof("[pkg/postgres] успешно завершены миграции к базе данных: postgres://%s:***@%s:%s/%s?sslmode=disable",
				cfg.Username, cfg.Host, cfg.Port, cfg.Database)
		}
	}()

	if err := m.Up(); err != nil {

		if err == migrate.ErrNoChange {
			v, _, _ := m.Version()
			logrus.Infof("[pkg/postgres] нет изменений при выполнении SQL миграции (текущая версия: %d)", v)
			return nil
		}

		return fmt.Errorf("migration pending error: %v", err)
	}

	return nil
}

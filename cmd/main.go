package main

import (
	"context"
	"time"

	"github.com/sater-151/todo-list/internal/controller/rest"
	settings "github.com/sater-151/todo-list/internal/init"
	pgRepo "github.com/sater-151/todo-list/internal/repository/postgres"
	"github.com/sater-151/todo-list/internal/service"
	"github.com/sater-151/todo-list/pkg/postgres"
	"github.com/sater-151/todo-list/pkg/server"
	"github.com/sirupsen/logrus"
)

const version = "v1.0.0"

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC822,
	})
}

func main() {
	appSettings, err := settings.Init()
	if err != nil {
		logrus.Fatalf("Ошибка при инициализации приложения: %s", err)
	}

	logrus.Infof("Запуск приложения %s", version)

	if appSettings.Arguments.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	pgConfig := postgres.StorageConfig{
		Database:        appSettings.Config.Database.DbName,
		Host:            appSettings.Config.Database.Host,
		Port:            appSettings.Config.Database.Port,
		Username:        appSettings.Config.Database.User,
		Password:        appSettings.Config.Database.Password,
		ApplicationName: "todo-list/" + version,
	}

	pgConn, err := postgres.NewPostgres(context.Background(), 3, 3*time.Second, pgConfig)
	if err != nil {
		logrus.Fatalf("Ошибка при инициализации базы данных: %s", err)
	}

	if appSettings.Arguments.Migrate {
		migrateConfig := postgres.MigrationConfig{
			MigrationPath: "file://./migrations",
			MigrationName: "migrations-todo-list",
		}

		if err := postgres.BeginMigrations(pgConfig, migrateConfig); err != nil {
			logrus.Fatalf("Ошибка при миграции базы данных: %s", err)
		}
	}

	postgresRepo, err := pgRepo.NewRepository(pgConn, appSettings.Config.Database.Schema, appSettings.Config.Database.CryproKey)
	if err != nil {
		logrus.Fatalf("Ошибка при инициализации репозитория: %s", err)
	}

	todoTask := service.New(postgresRepo, appSettings.Config.SecretKey)

	rst := rest.New(todoTask)

	srv := new(server.Server)

	if err := srv.Run(appSettings.Parameters.Server.Port, rst.Run()); err != nil {
		logrus.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}

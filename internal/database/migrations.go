package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func ApplyMigrations(dbURL string) {
	m, err := migrate.New("file://internal//database//schema", dbURL)
	if err != nil {
		logrus.Fatalf("failed to create migrator: %s", err.Error())
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("failed to apply migrations: %s", err.Error())
	}

	logrus.Print("Migrations applied successfully")
}

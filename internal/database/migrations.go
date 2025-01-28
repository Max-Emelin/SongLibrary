package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func ApplyMigrations(dbURL string) {
	logrus.Debug("Starting migration process...")

	m, err := migrate.New("file://internal/database/schema", dbURL)
	if err != nil {
		logrus.Errorf("Failed to create migrator: %v", err)
		logrus.Fatal(err)
	}
	logrus.Debug("Migrator created successfully.")

	logrus.Debug("Applying migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Errorf("Failed to apply migrations: %v", err)
		logrus.Fatal(err)
	}
}

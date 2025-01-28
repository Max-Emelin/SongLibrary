package database

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func ConnectDB(dbURL string) (*sqlx.DB, error) {
	logrus.Debug("Connecting to database...")
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		logrus.Errorf("Failed to initialize db: %s", err.Error())
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	return db, nil
}

func CreateDatabaseIfNotExists(dbURL, dbName string) error {
	logrus.Debug("Checking if database exists...")
	tempDBURL := dbURL[:strings.LastIndex(dbURL, "/")] + "?sslmode=disable"
	tempDB, err := sqlx.Open("postgres", tempDBURL)
	if err != nil {
		logrus.Errorf("Failed to connect to database server: %s", err.Error())
		return fmt.Errorf("failed to connect to database server: %w", err)
	}
	defer tempDB.Close()

	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		logrus.Errorf("Failed to create database: %s", err.Error())
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

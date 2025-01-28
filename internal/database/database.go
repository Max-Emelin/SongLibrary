package database

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func ConnectDB(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}
	return db, nil
}

func CreateDatabaseIfNotExists(dbURL, dbName string) error {
	tempDBURL := dbURL[:strings.LastIndex(dbURL, "/")] + "?sslmode=disable"
	tempDB, err := sqlx.Open("postgres", tempDBURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database server: %w", err)
	}
	defer tempDB.Close()

	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to create database: %w", err)
	}

	logrus.Printf("Database %s created or already exists", dbName)
	return nil
}

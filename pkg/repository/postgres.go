package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	songsTable  = "songs"
	lyricsTable = "lyrics"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	logrus.Debug("NewPostgresDB - creating connection string")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
	logrus.Debug("NewPostgresDB - connection string created")

	logrus.Debug("NewPostgresDB - opening database connection")
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		logrus.Error("NewPostgresDB - failed to open database connection", err)
		return nil, err
	}

	logrus.Debug("NewPostgresDB - pinging database")
	err = db.Ping()
	if err != nil {
		logrus.Error("NewPostgresDB - failed to ping database", err)
		return nil, err
	}

	logrus.Info("NewPostgresDB - database connection established successfully")
	return db, nil
}

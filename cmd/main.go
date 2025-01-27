package main

import (
	songlibrary "SongLibrary"
	"SongLibrary/pkg/apiClient"
	"SongLibrary/pkg/handler"
	"SongLibrary/pkg/repository"
	"SongLibrary/pkg/service"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	dbURL := "postgres://" + viper.GetString("db.username") + ":" + os.Getenv("DB_PASSWORD") + "@" +
		viper.GetString("db.host") + ":" + viper.GetString("db.port") + "/" + viper.GetString("db.dbname") + "?sslmode=" + viper.GetString("db.sslmode")

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	if err := createDatabaseIfNotExists(dbURL, viper.GetString("db.dbname")); err != nil {
		logrus.Fatalf("error creating database: %s", err.Error())
	}

	applyMigrations(dbURL)

	client := apiClient.NewClient("http://localhost:8080")
	repos := repository.NewRepository(db)
	service := service.NewService(repos, client)
	handlers := handler.NewHandler(service)

	srv := new(songlibrary.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print("SongLibrary Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("SongLibrary Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func applyMigrations(dbURL string) {
	m, err := migrate.New(
		"file://schema",
		dbURL,
	)
	if err != nil {
		logrus.Fatalf("failed to create migrator: %s", err.Error())
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("failed to apply migrations: %s", err.Error())
	}

	logrus.Print("Migrations applied successfully")
}

func createDatabaseIfNotExists(dbURL, dbName string) error {
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

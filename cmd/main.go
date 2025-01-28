package main

import (
	"SongLibrary/internal/database"
	"SongLibrary/internal/server"
	"SongLibrary/pkg/apiClient"
	"SongLibrary/pkg/handler"
	"SongLibrary/pkg/repository"
	"SongLibrary/pkg/service"
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	logrus.Info("Loading environment variables...")

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables: %s", err.Error())
	}
	logrus.Info("Environment variables loaded.")

	dbURL := getDatabaseURL()

	logrus.Info("Connecting to the database...")
	dbConn, err := database.ConnectDB(dbURL)
	if err != nil {
		logrus.Fatalf("Failed to connect to DB: %s", err.Error())
	}
	logrus.Info("Connected to the database successfully.")

	logrus.Info("Checking if the database exists...")
	if err := database.CreateDatabaseIfNotExists(dbURL, "song_library"); err != nil {
		logrus.Fatalf("Error creating database: %s", err.Error())
	}
	logrus.Info("Database created or already exists.")

	logrus.Info("Applying migrations...")
	database.ApplyMigrations(dbURL)
	logrus.Info("Migrations applied successfully.")

	logrus.Info("Initializing API client...")
	client := apiClient.NewClient("http://localhost:8080")
	logrus.Info("API client initialized.")

	repos := repository.NewRepository(dbConn)
	service := service.NewService(repos, client)
	handlers := handler.NewHandler(service)

	srv := new(server.Server)

	logrus.Info("Starting the HTTP server...")
	go func() {
		if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occurred while running HTTP server: %s", err.Error())
		}
	}()
	logrus.Info("HTTP server started.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Info("SongLibrary is shutting down...")

	logrus.Info("Shutting down the server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occurred on server shutdown: %s", err.Error())
	}

	logrus.Info("Closing the database connection...")
	if err := dbConn.Close(); err != nil {
		logrus.Errorf("Error occurred while closing the DB connection: %s", err.Error())
	}
	logrus.Info("Database connection closed.")
}

func getDatabaseURL() string {
	logrus.Debug("Building database URL...")
	url := "postgres://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") +
		"?sslmode=" + os.Getenv("DB_SSLMODE")
	logrus.Debugf("Database URL: %s", url)

	return url
}

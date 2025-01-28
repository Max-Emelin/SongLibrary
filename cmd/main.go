package main

import (
	"SongLibrary/internal/configs"
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
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	dbURL := getDatabaseURL()

	dbConn, err := database.ConnectDB(dbURL)
	if err != nil {
		logrus.Fatalf("failed to connect to DB: %s", err.Error())
	}

	if err := database.CreateDatabaseIfNotExists(dbURL, "song_library"); err != nil {
		logrus.Fatalf("error creating database: %s", err.Error())
	}

	database.ApplyMigrations(dbURL)

	client := apiClient.NewClient("http://localhost:8080")
	repos := repository.NewRepository(dbConn)
	service := service.NewService(repos, client)
	handlers := handler.NewHandler(service)

	srv := new(server.Server)
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

	if err := dbConn.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func getDatabaseURL() string {
	return "postgres://" + viper.GetString("db.username") + ":" + os.Getenv("DB_PASSWORD") + "@" +
		viper.GetString("db.host") + ":" + viper.GetString("db.port") + "/" + viper.GetString("db.dbname") +
		"?sslmode=" + viper.GetString("db.sslmode")
}

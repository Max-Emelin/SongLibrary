package configs

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config")

	logrus.Debug("Attempting to load the configuration...")

	if err := viper.ReadInConfig(); err != nil {
		dir, _ := os.Getwd()
		logrus.Errorf("Error reading config file. Current working directory: %s", dir)
		return err
	}

	logrus.Info("Configuration loaded successfully")
	return nil
}

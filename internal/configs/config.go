package configs

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		dir, _ := os.Getwd()
		logrus.Errorf("Current working directory: %s", dir)
		return err
	}

	logrus.Print("Configuration loaded successfully")
	return nil
}

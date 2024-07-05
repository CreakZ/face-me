package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	RedisHost     = "REDIS_HOST"
	RedisPort     = "REDIS_PORT"
	RedisPassword = "REDIS_PASSWORD"
)

func InitConfig() error {
	path, _ := os.Getwd()
	path = filepath.Join("../deploy")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AddConfigPath(path)
	// .env file can be located in root
	viper.AddConfigPath("")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

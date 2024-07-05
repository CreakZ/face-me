package database

import (
	"faceit_parser/pkg/config"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", viper.GetString(config.RedisHost), viper.GetInt(config.RedisPort)),
		Password: viper.GetString(config.RedisPassword),
		DB:       0,
	})

	return db
}

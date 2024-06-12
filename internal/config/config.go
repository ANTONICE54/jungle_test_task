package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBHost        string `mapstructure:"DB_HOST"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
}

// Loads config from .env file
func LoadConfig() Config {
	cfg := Config{}

	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return cfg

}

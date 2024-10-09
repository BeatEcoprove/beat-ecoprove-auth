package config

import "github.com/spf13/viper"

var config *Config

type Config struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
}

func LoadEnv(path string) {
	viper.SetConfigName(path)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		viper.Reset()
		viper.AutomaticEnv()
	}

	config = &Config{
		POSTGRES_PASSWORD: viper.GetString("POSTGRES_PASSWORD"),
		POSTGRES_USER:     viper.GetString("POSTGRES_USER"),
		POSTGRES_DB:       viper.GetString("POSTGRES_DB"),
		POSTGRES_HOST:     viper.GetString("POSTGRES_HOST"),
		POSTGRES_PORT:     viper.GetString("POSTGRES_PORT"),
	}
}

func GetCofig() *Config {
	return config
}

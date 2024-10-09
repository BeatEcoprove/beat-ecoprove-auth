package config

import "github.com/spf13/viper"

var config *Config

type Config struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     string

	BEAT_IDENTITY_SERVER uint16
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

		BEAT_IDENTITY_SERVER: viper.GetUint16("BEAT_IDENTITY_SERVER"),
	}
}

func GetCofig() *Config {
	return config
}

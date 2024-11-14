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

	JWT_AUDIENCE        string
	JWT_ISSUER          string
	JWT_ACCESS_EXPIRED  int
	JWT_REFRESH_EXPIRED int
	JWT_SECRET          string

	REDIS_HOST string
	REDIS_PORT string
	REDIS_DB   int
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

		JWT_AUDIENCE:        viper.GetString("JWT_AUDIENCE"),
		JWT_ISSUER:          viper.GetString("JWT_ISSUER"),
		JWT_ACCESS_EXPIRED:  viper.GetInt("JWT_ACCESS_EXPIRED"),
		JWT_REFRESH_EXPIRED: viper.GetInt("JWT_REFRESH_EXPIRED"),
		JWT_SECRET:          viper.GetString("JWT_SECRET"),

		REDIS_HOST: viper.GetString("REDIS_HOST"),
		REDIS_PORT: viper.GetString("REDIS_PORT"),
		REDIS_DB:   viper.GetInt("REDIS_DB"),
	}
}

func GetCofig() *Config {
	return config
}

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

	RABBIT_MQ_HOST           string
	RABBIT_MQ_PORT           string
	RABBITMQ_DEFAULT_USER    string
	RABBITMQ_DEFAULT_PASS    string
	RABBITMQ_DEFAULT_VHOST   string
	RABBIT_MQ_EXCHANGE       string
	RABBIT_MQ_QUEUE_ENTITIES string
	RABBIT_MQ_QUEUE_MAIL     string
	RABBIT_MQ_ROUTING_KEY    string
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

		RABBIT_MQ_HOST:           viper.GetString("RABBIT_MQ_HOST"),
		RABBIT_MQ_PORT:           viper.GetString("RABBIT_MQ_PORT"),
		RABBITMQ_DEFAULT_USER:    viper.GetString("RABBITMQ_DEFAULT_USER"),
		RABBITMQ_DEFAULT_PASS:    viper.GetString("RABBITMQ_DEFAULT_PASS"),
		RABBITMQ_DEFAULT_VHOST:   viper.GetString("RABBITMQ_DEFAULT_VHOST"),
		RABBIT_MQ_EXCHANGE:       viper.GetString("RABBIT_MQ_EXCHANGE"),
		RABBIT_MQ_QUEUE_ENTITIES: viper.GetString("RABBIT_MQ_QUEUE_ENTITIES"),
		RABBIT_MQ_QUEUE_MAIL:     viper.GetString("RABBIT_MQ_QUEUE_MAIL"),
		RABBIT_MQ_ROUTING_KEY:    viper.GetString("RABBIT_MQ_ROUTING_KEY"),
	}
}

func GetCofig() *Config {
	return config
}

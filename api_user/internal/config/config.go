package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/reversersed/LitGO-backend/tree/main/api_user/pkg/mongo"
	"github.com/reversersed/LitGO-backend/tree/main/api_user/pkg/rabbitmq"
)

type Config struct {
	Server   *ServerConfig
	Database *mongo.DatabaseConfig
	Rabbit   *rabbitmq.RabbitConfig
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-required:"true" env-description:"Server listening address"`
	Port        int    `env:"SERVER_PORT" env-required:"true" env-description:"Server listening port"`
	Environment string `env:"ENVIRONMENT" env-default:"debug" env-description:"Application environment"`
	JwtSecret   string `env:"JWT_SECRET" env-required:"true"  env-description:"JWT secret token. Must be unique and strong"`
}

var once sync.Once
var config *Config

func GetConfig() (*Config, error) {
	var e error
	once.Do(func() {
		server := &ServerConfig{}
		database := &mongo.DatabaseConfig{}
		rabbit := &rabbitmq.RabbitConfig{}

		if err := cleanenv.ReadConfig("config/.env", server); err != nil {
			var header string = "Server part config"
			desc, _ := cleanenv.GetDescription(server, &header)
			e = fmt.Errorf("%w: %s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", database); err != nil {
			var header string = "Database part config"
			desc, _ := cleanenv.GetDescription(database, &header)
			e = fmt.Errorf("%w\n%s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", rabbit); err != nil {
			var header string = "RabbitMQ part config"
			desc, _ := cleanenv.GetDescription(database, &header)
			e = fmt.Errorf("%w\n%s", err, desc)
			return
		}
		config = &Config{
			Server:   server,
			Database: database,
			Rabbit:   rabbit,
		}
	})
	if e != nil {
		return nil, e
	}
	return config, nil
}

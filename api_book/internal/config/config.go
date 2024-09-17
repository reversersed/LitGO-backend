package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/mongo"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/rabbitmq"
)

type Config struct {
	Server   *ServerConfig
	Database *mongo.DatabaseConfig
	Rabbit   *rabbitmq.RabbitConfig
}

type ServerConfig struct {
	Host             string `env:"SERVER_HOST" env-required:"true" env-description:"Server listening address"`
	Port             int    `env:"SERVER_PORT" env-required:"true" env-description:"Server listening port"`
	Environment      string `env:"ENVIRONMENT" env-default:"debug" env-description:"Application environment"`
	GenreServiceUrl  string `env:"SERVICE_GENRE_URL" env-required:"true" env-description:"External URL of genre service"`
	AuthorServiceUrl string `env:"SERVICE_AUTHOR_URL" env-required:"true" env-description:"External URL of author service"`
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
			e = fmt.Errorf("%v: %s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", database); err != nil {
			var header string = "Database part config"
			desc, _ := cleanenv.GetDescription(database, &header)
			e = fmt.Errorf("%v\n%s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", rabbit); err != nil {
			var header string = "RabbitMQ part config"
			desc, _ := cleanenv.GetDescription(database, &header)
			e = fmt.Errorf("%v\n%s", err, desc)
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

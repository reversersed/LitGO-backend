package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/reversersed/LitGO-backend-pkg/rabbitmq"
)

type Config struct {
	Server *ServerConfig
	Rabbit *rabbitmq.RabbitConfig
	File   *FileConfig
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-required:"true" env-description:"Server listening address"`
	Port        int    `env:"SERVER_PORT" env-required:"true" env-description:"Server listening port"`
	Environment string `env:"ENVIRONMENT" env-default:"debug" env-description:"Application environment"`
}

type FileConfig struct {
	BooksFolder        string `env:"FOLDER_BOOK_FILE" env-required:"true" env-description:"Book's epub folder"`
	BookCoversFolder   string `env:"FOLDER_COVERS_FILE" env-required:"true" env-description:"Book cover's folder"`
	AuthorCoversFolder string `env:"FOLDER_COVERS_AUTHOR" env-required:"true" env-description:"Authors cover's folder"`
}

var once sync.Once
var config *Config

func GetConfig() (*Config, error) {
	var e error
	once.Do(func() {
		server := &ServerConfig{}
		rabbit := &rabbitmq.RabbitConfig{}
		file := &FileConfig{}

		if err := cleanenv.ReadConfig("config/.env", server); err != nil {
			var header string = "Server part config"
			desc, _ := cleanenv.GetDescription(server, &header)
			e = fmt.Errorf("%v: %s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", rabbit); err != nil {
			var header string = "RabbitMQ part config"
			desc, _ := cleanenv.GetDescription(rabbit, &header)
			e = fmt.Errorf("%v\n%s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", file); err != nil {
			var header string = "File folder's part config"
			desc, _ := cleanenv.GetDescription(file, &header)
			e = fmt.Errorf("%v\n%s", err, desc)
			return
		}
		config = &Config{
			Server: server,
			Rabbit: rabbit,
			File:   file,
		}
	})
	if e != nil {
		return nil, e
	}
	return config, nil
}

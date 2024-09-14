package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server *ServerConfig
	Url    *UrlConfig
}
type UrlConfig struct {
	UserServiceUrl   string `env:"SERVICE_USER_URL" env-required:"true" env-description:"External URL of user (identity) service"`
	GenreServiceUrl  string `env:"SERVICE_GENRE_URL" env-required:"true" env-description:"External URL of genre service"`
	AuthorServiceUrl string `env:"SERVICE_AUTHOR_URL" env-required:"true" env-description:"External URL of author service"`
	BookServiceUrl   string `env:"SERVICE_BOOK_URL" env-required:"true" env-description:"External URL of book service"`
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
		url := &UrlConfig{}

		if err := cleanenv.ReadConfig("config/.env", server); err != nil {
			desc, _ := cleanenv.GetDescription(server, nil)
			e = fmt.Errorf("%w: %s", err, desc)
			return
		}
		if err := cleanenv.ReadConfig("config/.env", url); err != nil {
			desc, _ := cleanenv.GetDescription(url, nil)
			e = fmt.Errorf("%w: %s", err, desc)
			return
		}
		config = &Config{
			Server: server,
			Url:    url,
		}
	})
	if e != nil {
		return nil, e
	}
	return config, nil
}

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
	UserServiceUrl string `env:"SERVICE_USER_URL" env-required:"true"`
}
type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-required:"true"`
	Port        int    `env:"SERVER_PORT" env-required:"true"`
	JwtSecret   string `env:"JWT_SECRET" env-required:"true"`
	Environment string `env:"ENVIRONMENT"`
}

var once sync.Once
var config *Config

func GetConfig() (*Config, error) {
	var e error
	once.Do(func() {
		server := &ServerConfig{}
		url := &UrlConfig{}

		if err := cleanenv.ReadConfig("config/.env", server); err != nil {
			desc, _ := cleanenv.GetDescription(config, nil)
			e = fmt.Errorf("%v: %s", err, desc)
			return
		}
		if len(server.Environment) == 0 {
			server.Environment = "debug"
		}
		if err := cleanenv.ReadConfig("config/.env", url); err != nil {
			desc, _ := cleanenv.GetDescription(config, nil)
			e = fmt.Errorf("%v: %s", err, desc)
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

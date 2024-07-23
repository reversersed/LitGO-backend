package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
}
type DatabaseConfig struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_NAME" env-required:"true"`
	Password string `env:"DB_PASS" env-required:"true"`
	Base     string `env:"DB_BASE" env-required:"true"`
	AuthDb   string `env:"DB_AUTHDB" env-required:"true"`
}
type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-required:"true"`
	Port        int    `env:"SERVER_PORT" env-required:"true"`
	Environment string `env:"ENVIRONMENT"`
}

var once sync.Once
var config *Config

func GetConfig() (*Config, error) {
	var e error
	once.Do(func() {
		server := &ServerConfig{}
		database := &DatabaseConfig{}

		if err := cleanenv.ReadConfig("config/.env", server); err != nil {
			desc, _ := cleanenv.GetDescription(config, nil)
			e = fmt.Errorf("%v: %s", err, desc)
			return
		}
		if len(server.Environment) == 0 {
			server.Environment = "debug"
		}
		if err := cleanenv.ReadConfig("config/.env", database); err != nil {
			desc, _ := cleanenv.GetDescription(config, nil)
			e = fmt.Errorf("%v: %s", err, desc)
			return
		}
		config = &Config{
			Server:   server,
			Database: database,
		}
	})
	if e != nil {
		return nil, e
	}
	return config, nil
}

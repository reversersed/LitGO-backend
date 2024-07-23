package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/reversersed/go-grpc/tree/main/api_user/pkg/mongo"
)

type Config struct {
	Server   *ServerConfig
	Database *mongo.DatabaseConfig
}

type ServerConfig struct {
	Host        string `env:"SERVER_HOST" env-required:"true"`
	Port        int    `env:"SERVER_PORT" env-required:"true"`
	Environment string `env:"ENVIRONMENT"`
	JwtSecret   string `env:"JWT_SECRET" env-required:"true"`
}

var once sync.Once
var config *Config

func GetConfig() (*Config, error) {
	var e error
	once.Do(func() {
		server := &ServerConfig{}
		database := &mongo.DatabaseConfig{}

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

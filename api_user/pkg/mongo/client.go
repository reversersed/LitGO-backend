package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_NAME" env-required:"true"`
	Password string `env:"DB_PASS" env-required:"true"`
	Base     string `env:"DB_BASE" env-required:"true"`
	AuthDb   string `env:"DB_AUTHDB" env-required:"true"`
}

func NewClient(ctx context.Context, cfg *DatabaseConfig) (*mongo.Database, error) {
	var mongoURL string
	var anonymous bool

	if cfg.User == "" || cfg.Password == "" {
		anonymous = true
		mongoURL = fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	} else {
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	}
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(mongoURL)
	if !anonymous {
		clientOptions.SetAuth(options.Credential{
			Username:    cfg.User,
			Password:    cfg.Password,
			PasswordSet: true,
			AuthSource:  cfg.AuthDb,
		})
	}
	client, err := mongo.Connect(reqCtx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	return client.Database(cfg.Base), nil
}

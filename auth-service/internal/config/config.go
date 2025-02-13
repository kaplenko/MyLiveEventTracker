package config

import (
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	JWTSecret string        `env:"JWT_SECRET" required:"true"`
	TokenTTL  time.Duration `env:"TOKEN_TTL" default:"5h"`

	DBDSN string `env:"DB_DSN" required:"true"`

	AppHost string `env:"APP_HOST" required:"true"`
	AppPort string `env:"APP_PORT" required:"true"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	return cfg
}

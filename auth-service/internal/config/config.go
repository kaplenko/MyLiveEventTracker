package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"time"
)

type Config struct {
	JWTSecret string        `env:"JWT_SECRET" required:"true"`
	TokenTTL  time.Duration `env:"TOKEN_TTL" default:"5h"`

	DBUser     string `env:"DB_USER" required:"true"`
	DBPassword string `env:"DB_PASSWORD" required:"true"`
	DBName     string `env:"DB_NAME" required:"true"`
	DBHost     string `env:"DB_HOST" required:"true"`
	DBPort     string `env:"DB_PORT" required:"true"`

	AppHost string `env:"APP_HOST" required:"true"`
	AppPort string `env:"APP_PORT" required:"true"`

	GithubClientID     string `env:"GITHUB_CLIENT_ID" required:"true"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET" required:"true"`
	GithubRedirectURL  string `env:"GITHUB_REDIRECT_URL" required:"true"`
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка загрузки .env:", err)
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	return cfg
}

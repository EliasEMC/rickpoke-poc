package config

import (
    "github.com/caarlos0/env/v6"
)

type DatabaseConfig struct {
    Type     string `env:"DB_TYPE" envDefault:"postgres"`
    Host     string `env:"DB_HOST" envDefault:"localhost"`
    Port     int    `env:"DB_PORT" envDefault:"5432"`
    User     string `env:"DB_USER" envDefault:"alaska-eng"`
    Password string `env:"DB_PASSWORD"`
    Name     string `env:"DB_NAME" envDefault:"rickpoke"`
}

func NewDatabaseConfig() *DatabaseConfig {
    cfg := &DatabaseConfig{}
    if err := env.Parse(cfg); err != nil {
        panic(err)
    }
    return cfg
} 
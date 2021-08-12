package config

import "github.com/caarlos0/env"

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	GrpcPort int `env:"PORT" envDefault:"9092"`
	EnableReflection bool `env:"ENABLE_REFLECTION" envDefault:"false"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	return cfg, env.Parse(cfg)
}
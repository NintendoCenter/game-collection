package config

import "github.com/caarlos0/env"

type Config struct {
	LogLevel         string `env:"LOG_LEVEL" envDefault:"debug"`
	GrpcPort         int    `env:"PORT" envDefault:"9092"`
	EnableReflection bool   `env:"ENABLE_REFLECTION" envDefault:"false"`
	MongoConnection  string `env:"DATABASE_URL" envDefault:"mongodb://localhost:27017/nintendo-center"`
	GamesTopic       string `env:"GAMES_TOPIC" envDefault:"sync_games"`
	QueueAddr        string `env:"QUEUE_ADDR" envDefault:"localhost:4150"`
	ElasticAdds      string `env:"ELASTIC_ADDR" envDefault:"http://localhost:9200"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	return cfg, env.Parse(cfg)
}
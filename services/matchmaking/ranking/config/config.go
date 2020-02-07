package config

import (
	"github.com/4726/game/pkg/config"
)

type Config struct {
	Redis                RedisConfig
	NSQ               NSQConfig
	Port              int
	Metrics           MetricsConfig
}

type MetricsConfig struct {
	Port  int
	Route string
}

type RedisConfig struct {
	Addr string
	Password string
	DB int
	SetName string
}

type NSQConfig struct {
	Addr, Topic, Channel string
}

const defaultMaxMatchResponses = 100
const defaultPort = 14000

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "history",
		Defaults: map[string]interface{}{
			"Port":              defaultPort,
			"Metrics.Port":      14001,
			"Metrics.Route":     "/metrics",
		},
		FilePath: filePath,
	})
	return cfg, err
}

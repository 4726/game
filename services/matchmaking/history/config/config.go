package config

import (
	"github.com/4726/game/pkg/config"
)

type Config struct {
	DB                DBConfig
	NSQ               NSQConfig
	MaxMatchResponses uint32
	Port              int
	Metrics           MetricsConfig
}

type MetricsConfig struct {
	Port  int
	Route string
}

type DBConfig struct {
	Name, Collection string
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
			"MaxMatchResponses": defaultMaxMatchResponses,
			"Port":              defaultPort,
			"Metrics.Port":      14001,
			"Metrics.Route":     "/metrics",
		},
		FilePath: filePath,
	})
	return cfg, err
}

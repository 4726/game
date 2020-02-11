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
	TLS               TLSConfig
}

type MetricsConfig struct {
	Port  int
	Route string
}

type DBConfig struct {
	Name, Collection string
	Addr             string
	DialTimeout      uint //seconds
}

type NSQConfig struct {
	Addr, Topic, Channel string
}

type TLSConfig struct {
	CertPath, KeyPath string
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
			"DB.Name":           "matchmaking",
			"DB.Collection":     "history",
			"DB.Addr":           "mongodb://localhost:27017",
			"DB.DialTimeout":    10,
		},
		FilePath: filePath,
	})
	return cfg, err
}

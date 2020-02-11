package config

import (
	"github.com/4726/game/pkg/config"
)

type Config struct {
	DB      DBConfig
	Port    int
	Metrics MetricsConfig
	TLS     TLSConfig
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

type TLSConfig struct {
	CertPath, KeyPath string
}

const defaultMaxMatchResponses = 100
const defaultPort = 14000

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "live",
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

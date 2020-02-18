package config

import (
	"github.com/4726/game/pkg/config"
)

type Config struct {
	Port    int
	Metrics MetricsConfig
	TLS     TLSConfig
}

type MetricsConfig struct {
	Port  int
	Route string
}

type TLSConfig struct {
	CertPath, KeyPath string
}

const defaultMaxMatchResponses = 100
const defaultPort = 14000

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "custommatch",
		Defaults: map[string]interface{}{
			"Port":              defaultPort,
			"Metrics.Port":      14001,
			"Metrics.Route":     "/metrics",
		},
		FilePath: filePath,
	})
	return cfg, err
}

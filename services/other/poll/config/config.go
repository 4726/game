package config

import (
	"github.com/4726/game/pkg/config"
)

type Config struct {
	Redis   RedisConfig
	Port    int
	Metrics MetricsConfig
	TLS     TLSConfig
	MaxPollChoices uint
	MaxExpireMinutes uint
}

type MetricsConfig struct {
	Port  int
	Route string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type TLSConfig struct {
	CertPath, KeyPath string
}

const defaultMaxMatchResponses = 100
const defaultPort = 14000
const defaultMaxPollChoices = 10
const defaultMaxExpireMinutes = 144000

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "poll",
		Defaults: map[string]interface{}{
			"Port":          defaultPort,
			"Metrics.Port":  14001,
			"Metrics.Route": "/metrics",
			"MaxPollChoices": defaultMaxPollChoices,
			"MaxExpireMinutes": defaultMaxExpireMinutes,
		},
		FilePath: filePath,
	})
	return cfg, err
}

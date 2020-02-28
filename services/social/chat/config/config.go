package config

import (
	"github.com/4726/game/pkg/config"
)

type Config struct {
	Cassandra CassandraConfig
	Port      int
	Metrics   MetricsConfig
	TLS       TLSConfig
	NSQ       NSQConfig
}

type NSQConfig struct {
	Addr  string
	Topic string
}

type MetricsConfig struct {
	Port  int
	Route string
}

type CassandraConfig struct {
	Host        string
	Port        int
	DialTimeout uint
}

type TLSConfig struct {
	CertPath, KeyPath string
}

const defaultPort = 14000
const defaultCassandraPort = 9042

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "history",
		Defaults: map[string]interface{}{
			"Port":           defaultPort,
			"Metrics.Port":   14001,
			"Metrics.Route":  "/metrics",
			"Cassandra.Port": defaultCassandraPort,
		},
		FilePath: filePath,
	})
	return cfg, err
}

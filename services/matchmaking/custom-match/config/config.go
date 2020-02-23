package config

import (
	"github.com/4726/game/pkg/config"
)

type Config struct {
	Port    int
	Metrics MetricsConfig
	TLS     TLSConfig
	DB      DBConfig
}

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
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
const defaultDBName = "postgres"
const defaultDBUser = "postgres"
const defaultDBPassword = "postgres"
const defaultDBHost = "localhost"
const defaultDBPort = 5432

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "custommatch",
		Defaults: map[string]interface{}{
			"Port":          defaultPort,
			"Metrics.Port":  14001,
			"Metrics.Route": "/metrics",
			"DB.Name":       defaultDBName,
			"DB.User":       defaultDBUser,
			"DB.Password":   defaultDBPassword,
			"DB.Host":       defaultDBHost,
			"DB.Port":       defaultDBPort,
		},
		FilePath: filePath,
	})
	return cfg, err
}

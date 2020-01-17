package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/4726/game/pkg/config"
)

type Config struct {
	DB                DBConfig
	NSQ               NSQConfig
	MaxMatchResponses uint32
}

type DBConfig struct {
	Name, Collection string
}

type NSQConfig struct {
	Addr, Topic, Channel string
}

const defaultMaxMatchResponses = 100

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "history"
		Defaults: map[string]interface{}{"MaxMatchResponses", defaultMaxMatchResponses}
		FilePath: filePath,
	})
	return cfg, err
}
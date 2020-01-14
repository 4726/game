package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DB  DBConfig
	NSQ NSQConfig
}

type DBConfig struct {
	Name, Collection string
}

type NSQConfig struct {
	Addr, Topic, Channel string
}

func LoadConfig(filePath string) (Config, error) {
	if filePath != "" {
		ext := filepath.Ext(filePath)[1:]
		viper.SetConfigType(ext)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return Config{}, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
			return Config{}, err
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			return Config{}, err
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("history")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, v := range os.Environ() {
		tokens := strings.Split(v, "=")
		if !strings.HasPrefix(tokens[0], "HISTORY") {
			continue
		}
		key := strings.TrimPrefix(tokens[0], "HISTORY_")
		viper.BindEnv(key)
	}

	var cfg Config
	err := viper.Unmarshal(&cfg)
	return cfg, err
}

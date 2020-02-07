package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var defaultEnvKeyReplacer = strings.NewReplacer(".", "_")
var defaultConfigFileName = "config.yml"


type ConfigOpts struct {
	EnvPrefix      string
	EnvKeyReplacer *strings.Replacer
	Defaults       map[string]interface{}
	FilePath       string
}

func (o ConfigOpts) fix() (ConfigOpts, error) {
	if o.EnvKeyReplacer == nil {
		o.EnvKeyReplacer = defaultEnvKeyReplacer
	}
	if o.FilePath == "" {
		o.FilePath = defaultConfigFileName
	}
	if o.EnvPrefix == "" {
		return o, errors.New("EnvPrefix cannot be empty")
	}
	return o, nil
}

// LoadConfig loads configuration values into cfg
func LoadConfig(cfg interface{}, opts ConfigOpts) error {
	opts, err := opts.fix()
	if err != nil {
		return err
	}

	ext := filepath.Ext(opts.FilePath)[1:]
	viper.SetConfigType(ext)
	data, err := ioutil.ReadFile(opts.FilePath)
	if err != nil {
		return err
	}

	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return err
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(opts.EnvPrefix)
	opts.EnvPrefix = strings.ToUpper(opts.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for k, v := range opts.Defaults {
		viper.SetDefault(k, v)
	}
	for _, v := range os.Environ() {
		tokens := strings.Split(v, "=")
		if !strings.HasPrefix(tokens[0], opts.EnvPrefix) {
			continue
		}
		key := strings.TrimPrefix(tokens[0], fmt.Sprintf("%v_", opts.EnvPrefix))
		viper.BindEnv(key)
	}

	return viper.Unmarshal(cfg)
}

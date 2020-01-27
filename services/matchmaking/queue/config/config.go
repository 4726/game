package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/4726/game/pkg/config"
)

type Config struct {
	Limit int
	PerMatch int
	RatingRange int
	AcceptTimeoutSeconds int 
}

const defaultLimit = 10000
const defaultPerMatch = 10
const defaultRatingRange = 100
const defaultAcceptTimeoutSeconds = 20

func LoadConfig(filePath string) (Config, error) {
	var cfg Config
	err := config.LoadConfig(&cfg, config.ConfigOpts{
		EnvPrefix: "queue"
		Defaults: map[string]interface{}{
			"Limit": defaultLimit,
			"PerMatch": defaultPerMatch,
			"RatingRange": defaultRatingRange,
			"AcceptTimeoutSeconds": defaultAcceptTimeoutSeconds,
		}
		FilePath: filePath,
	})
	return cfg, err
}

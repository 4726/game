package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigDeafultFileName(t *testing.T) {
	cfg, err := LoadConfig("")
	assert.NoError(t, err)
	expectedCfg := Config{
		DBConfig{"dbname", "collectionname"},
		NSQConfig{"127.0.0.1:4150", "topic", "channel"},
		100,
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigCustomFileName(t *testing.T) {
	cfg, err := LoadConfig("config.yml")
	assert.NoError(t, err)
	expectedCfg := Config{
		DBConfig{"dbname", "collectionname"},
		NSQConfig{"127.0.0.1:4150", "topic", "channel"},
		100,
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigEnvVariablesDefaultFileName(t *testing.T) {
	os.Setenv("HISTORY_DB_NAME", "dbname2")
	os.Setenv("HISTORY_DB_COLLECTION", "collectionname2")
	os.Setenv("HISTORY_NSQ_ADDR", "127.0.0.1:4151")
	os.Setenv("HISTORY_NSQ_TOPIC", "topic2")
	os.Setenv("HISTORY_NSQ_CHANNEL", "channel2")
	os.Setenv("HISTORY_MAXMATCHRESPONSES", "102")
	cfg, err := LoadConfig("")
	assert.NoError(t, err)
	expectedCfg := Config{
		DBConfig{"dbname2", "collectionname2"},
		NSQConfig{"127.0.0.1:4151", "topic2", "channel2"},
		102,
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigEnvVariablesCustomFileName(t *testing.T) {
	os.Setenv("HISTORY_DB_NAME", "dbname2")
	os.Setenv("HISTORY_DB_COLLECTION", "collectionname2")
	os.Setenv("HISTORY_NSQ_ADDR", "127.0.0.1:4151")
	os.Setenv("HISTORY_NSQ_TOPIC", "topic2")
	os.Setenv("HISTORY_NSQ_CHANNEL", "channel2")
	os.Setenv("HISTORY_MAXMATCHRESPONSES", "102")
	cfg, err := LoadConfig("config.yml")
	assert.NoError(t, err)
	expectedCfg := Config{
		DBConfig{"dbname2", "collectionname2"},
		NSQConfig{"127.0.0.1:4151", "topic2", "channel2"},
		102,
	}
	assert.Equal(t, expectedCfg, cfg)
}

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	DB                testDBConfig
	NSQ               testNSQConfig
	MaxMatchResponses uint32 `mapstructure:"max_match_responses"`
	Empty             string
	Nested testNested
}

type testDBConfig struct {
	Name       string
	Collection string
}

type testNSQConfig struct {
	Addr, Topic, Channel string
}

type testNested struct {
	Number int
}

func TestLoadConfigDeafultFile(t *testing.T) {
	var cfg testConfig
	err := LoadConfig(&cfg, ConfigOpts{EnvPrefix: "history"})
	assert.NoError(t, err)
	expectedCfg := testConfig{
		testDBConfig{"dbname", "collectionname"},
		testNSQConfig{"127.0.0.1:4150", "topic", "channel"},
		100,
		"",
		testNested{0},
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigFile(t *testing.T) {
	var cfg testConfig
	err := LoadConfig(&cfg, ConfigOpts{EnvPrefix: "history", FilePath: "config.yml"})
	assert.NoError(t, err)
	expectedCfg := testConfig{
		testDBConfig{"dbname", "collectionname"},
		testNSQConfig{"127.0.0.1:4150", "topic", "channel"},
		100,
		"",
		testNested{0},
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigEnvPrefix(t *testing.T) {
	os.Setenv("HISTORY_DB_NAME", "dbname2")
	os.Setenv("HISTORY_DB_COLLECTION", "collectionname2")
	os.Setenv("HISTORY_NSQ_ADDR", "127.0.0.1:4151")
	os.Setenv("HISTORY_NSQ_TOPIC", "topic2")
	os.Setenv("HISTORY_NSQ_CHANNEL", "channel2")
	os.Setenv("HISTORY_MAX_MATCH_RESPONSES", "102")
	var cfg testConfig
	err := LoadConfig(&cfg, ConfigOpts{EnvPrefix: "history"})
	assert.NoError(t, err)
	expectedCfg := testConfig{
		testDBConfig{"dbname2", "collectionname2"},
		testNSQConfig{"127.0.0.1:4151", "topic2", "channel2"},
		102,
		"",
		testNested{0},
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigEnvPrefixWithFile(t *testing.T) {
	os.Setenv("HISTORY_DB_NAME", "dbname2")
	os.Setenv("HISTORY_DB_COLLECTION", "collectionname2")
	os.Setenv("HISTORY_NSQ_ADDR", "127.0.0.1:4151")
	os.Setenv("HISTORY_NSQ_TOPIC", "topic2")
	os.Setenv("HISTORY_NSQ_CHANNEL", "channel2")
	os.Setenv("HISTORY_MAXMATCHRESPONSES", "102")
	var cfg testConfig
	err := LoadConfig(&cfg, ConfigOpts{EnvPrefix: "history", FilePath: "config.yml"})
	assert.NoError(t, err)
	expectedCfg := testConfig{
		testDBConfig{"dbname2", "collectionname2"},
		testNSQConfig{"127.0.0.1:4151", "topic2", "channel2"},
		102,
		"",
		testNested{0},
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigDefault(t *testing.T) {
	os.Setenv("HISTORY_DB_NAME", "dbname2")
	os.Setenv("HISTORY_DB_COLLECTION", "collectionname2")
	os.Setenv("HISTORY_NSQ_ADDR", "127.0.0.1:4151")
	os.Setenv("HISTORY_NSQ_TOPIC", "topic2")
	os.Setenv("HISTORY_NSQ_CHANNEL", "channel2")
	os.Setenv("HISTORY_MAXMATCHRESPONSES", "102")
	var cfg testConfig
	defaults := map[string]interface{}{"Empty": "hello world"}
	err := LoadConfig(&cfg, ConfigOpts{EnvPrefix: "history", Defaults: defaults})
	assert.NoError(t, err)
	expectedCfg := testConfig{
		testDBConfig{"dbname2", "collectionname2"},
		testNSQConfig{"127.0.0.1:4151", "topic2", "channel2"},
		102,
		"hello world",
		testNested{0},
	}
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadConfigNestedDefault(t *testing.T) {
	os.Setenv("HISTORY_DB_NAME", "dbname2")
	os.Setenv("HISTORY_DB_COLLECTION", "collectionname2")
	os.Setenv("HISTORY_NSQ_ADDR", "127.0.0.1:4151")
	os.Setenv("HISTORY_NSQ_TOPIC", "topic2")
	os.Setenv("HISTORY_NSQ_CHANNEL", "channel2")
	os.Setenv("HISTORY_MAXMATCHRESPONSES", "102")
	var cfg testConfig
	defaults := map[string]interface{}{"nested.number": 5}
	err := LoadConfig(&cfg, ConfigOpts{EnvPrefix: "history", Defaults: defaults})
	assert.NoError(t, err)
	expectedCfg := testConfig{
		testDBConfig{"dbname2", "collectionname2"},
		testNSQConfig{"127.0.0.1:4151", "topic2", "channel2"},
		102,
		"hello world",
		testNested{5},
	}
	assert.Equal(t, expectedCfg, cfg)
}

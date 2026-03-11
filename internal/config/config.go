package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Redis  RedisConfig  `mapstructure:"redis"`
	Agent  AgentConfig  `mapstructure:"agent"`
	Memory MemoryConfig `mapstructure:"memory"`
}

type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AgentConfig struct {
	ID string `mapstructure:"id"`
}

type MemoryConfig struct {
	DefaultTTL int           `mapstructure:"default_ttl"`
	MaxResults int           `mapstructure:"max_results"`
	Embedding  EmbeddingConfig `mapstructure:"embedding"`
}

type EmbeddingConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Model   string `mapstructure:"model"`
}

var (
	cfg     *Config
	cfgFile string
)

func InitConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(filepath.Join(home, ".memorycli"))
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("MEMORYCLI")
	viper.AutomaticEnv()

	viper.SetDefault("redis.url", "redis://localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("agent.id", "")
	viper.SetDefault("memory.default_ttl", 0)
	viper.SetDefault("memory.max_results", 100)
	viper.SetDefault("memory.embedding.enabled", false)
	viper.SetDefault("memory.embedding.model", "text-embedding-ada-002")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func GetConfig() *Config {
	return cfg
}

func SetCfgFile(path string) {
	cfgFile = path
}

func GetRedisURL() string {
	return cfg.Redis.URL
}

func GetAgentID() string {
	if cfg.Agent.ID != "" {
		return cfg.Agent.ID
	}
	hostname, _ := os.Hostname()
	return hostname
}

package config

import (
	"flag"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	configVarName = "config"
)

// Config – ...
type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPCConfig    `yaml:"grpc"`
}

// GRPCConfig – ...
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// Load – ...
func Load() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("path to config is empty")
	}

	return parseConfig(path)
}

func parseConfig(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file doesn't exist")
	}

	var cfg Config

	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		panic("can't read config")
	}

	if err = v.Unmarshal(&cfg); err != nil {
		panic("can't unmarshal config")
	}

	return &cfg
}

func fetchConfigPath() (path string) {
	flag.StringVar(&path, configVarName, "", "path to config")
	flag.Parse()

	if path == "" {
		path = os.Getenv(configVarName)
	}

	return path
}

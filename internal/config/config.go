package config

import (
	"flag"
	"os"
	"time"

	"github.com/spf13/viper"

	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
)

const (
	configVarName = "config"
)

// Config – ...
type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	GRPC GRPCConfig `yaml:"grpc"`
}

// GRPCConfig – ...
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// Load – ...
func Load() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, errorpkg.New("path to config is empty")
	}

	return parseConfig(path)
}

func parseConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errorpkg.New("file doesn't exist")
	}

	var cfg Config

	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, errorpkg.New("can't read config")
	}

	if err = v.Unmarshal(&cfg); err != nil {
		return nil, errorpkg.New("can't unmarshal config")
	}

	return &cfg, nil
}

func fetchConfigPath() (path string) {
	flag.StringVar(&path, configVarName, "", "path to config")
	flag.Parse()

	if path == "" {
		path = os.Getenv(configVarName)
	}

	return path
}

package config

import (
	"flag"
	"fmt"
	"github.com/kalyuzhin/sso-service/internal/model"
	"net"
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
	Env      string     `yaml:"env" env-default:"local"`
	GRPC     GRPCConfig `yaml:"grpc"`
	Database DataBaseConfig
}

// GRPCConfig – ...
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// DataBaseConfig – ...
type DataBaseConfig struct {
	Database string `env:"DATABASE"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}

func (d *DataBaseConfig) GetDSN() string {
	switch d.Database {
	case model.PostgreSQLName:
		return fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable", d.User, d.Password,
			net.JoinHostPort(d.Host, d.Port), d.Name)
	default:
		return ""
	}
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

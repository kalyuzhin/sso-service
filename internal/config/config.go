package config

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"github.com/kalyuzhin/sso-service/internal/model"
)

const (
	configVarName = "config"
	envVarName    = "env-path"
)

// Config – ...
type Config struct {
	Env                    string        `yaml:"env" env-default:"local"`
	RefreshTokenExparation time.Duration `yaml:"refresh-token-exparation"`
	GRPC                   GRPCConfig    `yaml:"grpc"`
	Database               DataBaseConfig
}

// GRPCConfig – ...
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// DataBaseConfig – ...
type DataBaseConfig struct {
	Database string `env:"DATABASE"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USERNAME"`
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
	paths := fetchConfigPath(configVarName, envVarName)
	for _, path := range paths {
		if path == "" {
			return nil, errorpkg.New("path is empty")
		}
	}

	return parseConfig(paths[0], paths[1])
}

func parseConfig(pathConfig, pathEnv string) (*Config, error) {
	if _, err := os.Stat(pathConfig); os.IsNotExist(err) {
		return nil, errorpkg.New("file doesn't exist")
	}

	var cfg Config

	v := viper.New()
	v.SetConfigFile(pathConfig)
	err := v.ReadInConfig()
	if err != nil {
		return nil, errorpkg.New("can't read config")
	}

	if err = v.Unmarshal(&cfg); err != nil {
		return nil, errorpkg.New("can't unmarshal config")
	}

	if _, err = os.Stat(pathEnv); os.IsNotExist(err) {
		return nil, errorpkg.New("file doesn't exist")
	}

	var db DataBaseConfig
	if err = godotenv.Load(pathEnv); err != nil {
		return nil, errorpkg.WrapErr(err, "can't load env")
	}
	err = env.Parse(&db)
	if err != nil {
		return nil, errorpkg.WrapErr(err, "can't parse env")
	}

	cfg.Database = db

	return &cfg, nil
}

func fetchConfigPath(varNames ...string) []string {
	paths := make([]string, 2, 2)

	for idx, varName := range varNames {
		flag.StringVar(&paths[idx], varName, "", fmt.Sprintf("path to %s", varName))
		if paths[idx] == "" {
			paths[idx] = os.Getenv(varName)
		}
	}

	flag.Parse()

	return paths
}

package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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
	PrivateRSAKey          *rsa.PrivateKey
	PublicRSAKey           *rsa.PublicKey
	PathToRSAKey           string `env:"PATH_TO_PRIVATE_KEY"`
	PathToRSAPub           string `env:"PATH_TO_PUB_KEY"`
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

	cfg, err := parseConfig(paths[0], paths[1])
	if err != nil {
		return nil, err
	}

	err = cfg.readRSAPrivateKey()
	if err != nil {
		return nil, err
	}

	err = cfg.readRSAPublicKey()
	if err != nil {
		return nil, err
	}

	return cfg, nil
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

	if err = godotenv.Load(pathEnv); err != nil {
		return nil, errorpkg.WrapErr(err, "can't load env")
	}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, errorpkg.WrapErr(err, "can't parse env")
	}

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

func (c *Config) readRSAPrivateKey() error {
	bytes, err := os.ReadFile(c.PathToRSAKey)
	if err != nil {
		return errorpkg.WrapErr(err, "can't read private key file")
	}

	block, _ := pem.Decode(bytes)
	parseRes, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return errorpkg.WrapErr(err, "can't parse private key")
	}

	c.PrivateRSAKey = parseRes.(*rsa.PrivateKey)

	return nil
}

func (c *Config) readRSAPublicKey() error {
	bytes, err := os.ReadFile(c.PathToRSAPub)
	if err != nil {
		return errorpkg.WrapErr(err, "can't read private key file")
	}

	block, _ := pem.Decode(bytes)
	parseRes, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return errorpkg.WrapErr(err, "can't parse private key")
	}

	c.PublicRSAKey = parseRes.(*rsa.PublicKey)

	return nil
}

package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	Options  tokenOptions   `yaml:"token_options"`
	Database DatabaseConfig `yaml:"database" env-required:"true"`
	GRPC     GRPCConfig     `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type tokenOptions struct {
	JWTRefreshTTL      time.Duration `yaml:"token_refresh_ttl" env-required:"true"`
	JWTAccessTTL       time.Duration `yaml:"token_access_ttl" env-required:"true"`
	JWTVerificationTTL time.Duration `yaml:"token_verification_ttl" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg

}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parsed()

	if res == "" {
		res = os.Getenv("CONFIG_PATH_AUTH")
	}

	return res
}

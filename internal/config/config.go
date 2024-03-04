package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"grpc/internal/storage"
	"os"
	"time"
)

type Config struct {
	Env         string           `yaml:"env" env-default:"local"`
	GRPC        GRPCConfig       `yaml:"grpc"`
	Clients     ClientsConfig    `yaml:"clients"`
	Postgres    storage.Postgres `yaml:"postgres"`
	Coefficient float64          `yaml:"coefficient"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Client struct {
	Address      string        `yaml:"address"`
	RetriesCount int           `yaml:"retries_count"`
	Timeout      time.Duration `yaml:"timeout"`
	Insecure     bool          `yaml:"insecure"`
}

type ClientsConfig struct {
	Receiver Client `yaml:"receiver"`
}

func MustLoad() *Config {
	path, k := FetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	if k != 0 {
		cfg.Coefficient = k
	}

	return &cfg
}

func FetchConfigPath() (string, float64) {
	var path string
	var k float64

	// ./app --config=./path/to/config/file.yaml
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Float64Var(&k, "k", 0, "an STD anomaly coefficient")
	flag.Parse()

	if path == "" {
		// CONFIG_PATH=./path/to/config/file.yaml ./app
		path = os.Getenv("CONFIG_PATH")
	}

	return path, k
}

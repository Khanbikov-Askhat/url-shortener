package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address     string        `yaml:"address" default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" default:"60s"`
}

type Config struct {
	Env          string     `yaml:"env" default:"local"`
	StoragePath  string     `yaml:"storage_path" env-required:"true"`
	HTTPServer   HTTPServer `yaml:"http_server"`
	Port         string     `yaml:"port" default:"8080"`
	BaseURL      string     `yaml:"base_url" default:"http://localhost:8082"`
	DBPath       string     `yaml:"db_path" default:"./storage/storage.db"`
	LogLevel     string     `yaml:"log_level" default:"info"`
	LoggerFormat string     `yaml:"logger_format" default:"json"`
	LoggerOutput string     `yaml:"logger_output" default:"stdout"`
	LoggerLevel  string     `yaml:"logger_level" default:"info"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	//Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	//Load config
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg

}

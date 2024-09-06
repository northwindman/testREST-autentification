package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"20s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	// check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config File does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config file err: %s", err)
	}

	return &cfg
}

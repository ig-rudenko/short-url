package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	DSN        string `yaml:"storage_dsn" env-required:"true"`
	HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Load() *Config {
	config := os.Getenv("CONFIG")
	if config == "" {
		config = "./config/local.yaml"
	}

	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Fatalf("Файл конфигурации %s не существует", config)
	}

	var cgf Config

	if err := cleanenv.ReadConfig(config, &cgf); err != nil {
		log.Fatalf("Не удалось прочитать конфигурацию %s", err)
	}

	return &cgf
}

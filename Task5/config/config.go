package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		Comic  `yaml:"comic"`
		Server `yaml:"server"`
		DB     `yaml:"db"`
	}

	Comic struct {
		Source_url string `env-required:"true" yaml:"source_url"    env:"COMIC_SOURCE_URL"`
		Parallel   int    `env-required:"true" yaml:"parallel"    env:"COMIC_PARALLEL"`
	}

	Server struct {
		Port int `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	DB struct {
		DB_file string `env-required:"true" yaml:"db_file"   env:"DB_DB_FILE"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ReadConfig() *Config {
	// Configuration
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	return cfg
}

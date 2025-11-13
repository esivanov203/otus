package main

import (
	"fmt"
	internalhttp "github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/server/http"
	"os"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger  logger.Conf             `yaml:"logger"`
	Server  internalhttp.ServerConf `yaml:"server"`
	Storage StorageConf             `yaml:"storage"`
}

type StorageConf struct {
	Type string `yaml:"type"` // memory, sql
	Dsn  string `yaml:"dsn"`  // for sql
}

func NewConfig(configFile string) (Config, error) {
	cfg := Config{}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, err
	}
	expanded := os.ExpandEnv(string(data))

	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return cfg, fmt.Errorf("decoding %s: %w", configFile, err)
	}

	return cfg, nil
}

package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Server  ServerConf  `yaml:"server"`
	Storage StorageConf `yaml:"storage"`
}

type LoggerConf struct {
	Level string `yaml:"level"` // DEBUG, INFO, WARN, ERROR
	// TODO
}

type ServerConf struct {
	Host string `yaml:"host"` // listen interface
	Port int    `yaml:"port"` // listen port
}

type StorageConf struct {
	Type string `yaml:"type"` // memory, sql
	Dsn  string `yaml:"dsn"`  // for sql
}

func NewConfig(configFile string) (Config, error) {
	cfg := Config{}

	f, err := os.Open(configFile)
	if err != nil {
		return cfg, err
	}
	defer func() { _ = f.Close() }()

	d := yaml.NewDecoder(f)
	if err := d.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("decoding %s: %w", configFile, err)
	}

	return cfg, nil
}

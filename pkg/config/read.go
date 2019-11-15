package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Targets []string
	}

	ServiceConfig struct {
		Tasks []Task
	}

	Task struct {
		Name  string
		Files []File
	}

	File struct {
		Command  string
		Paths    []string
		Excluded bool
	}
)

func Read(cfgPath string, cfg *Config) error {
	f, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(cfg)
}

func ReadService(cfgPath string, cfg *ServiceConfig) error {
	f, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(cfg)
}

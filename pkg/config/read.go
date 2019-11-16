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
		Name   string
		Files  []File
		Change Change
	}

	File struct {
		Command  string
		Paths    []string
		Excluded bool
	}

	Change struct {
		IsFileChanged  IsFileChanged  `yaml:"is_file_changed"`
		IsFilesChanged IsFilesChanged `yaml:"is_files_changed"`
		ChangedFiles   ChangedFiles   `yaml:"changed_files"`
	}

	IsFileChanged struct {
		Command string
	}

	IsFilesChanged struct {
		Command string
		Stdin   bool
	}

	ChangedFiles struct {
		Command string
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

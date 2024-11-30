package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// Config is the configuration options for a screenshot.
type Config struct {
	Input    string
	Projects []Project

	// Settings
	Config      string
	Interactive bool

	Output         string
	Execute        string
	ExecuteTimeout time.Duration
}

var userConfigPath = filepath.Join(os.Getenv("O8ROOT"), "o8", "meta.ttl")

func loadUserConfig() (*Config, error) {
	f, readFileErr := loadUserConfigFile()

	c := &Config{}

	if readFileErr != nil {
		return nil, readFileErr
	}
	defer f.Close()

	byteValue, _ := io.ReadAll(f)

	marshallErr := json.Unmarshal(byteValue, c)

	if marshallErr != nil {
		return nil, marshallErr
	}

	return c, nil
}

func loadUserConfigFile() (fs.File, error) {
	return os.Open(userConfigPath)
}

// TODO: fix me by saving the ontology
func saveUserConfig(config Config) error {
	config.Input = ""
	config.Output = ""
	config.Interactive = false

	err := os.MkdirAll(filepath.Dir(userConfigPath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(userConfigPath)
	if err != nil {
		return err
	}
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = f.Write(b)

	print(userConfigPath, "SAVED")

	return err
}

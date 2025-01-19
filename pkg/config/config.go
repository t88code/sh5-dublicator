package config

import (
	"fmt"
	"gopkg.in/gcfg.v1"

	"os"
)

type Config struct {
	SH5SRC struct {
		BaseURL  string
		Username string
		Password string
		DebugLog bool
	}
	SH5DST struct {
		BaseURL  string
		Username string
		Password string
		DebugLog bool
	}
	SYNC struct {
		Refs string
	}
}

var cfg Config

func NewConfig() (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = gcfg.ReadFileInto(&cfg, fmt.Sprintf("%s/config.ini", pwd))
	if err != nil {
		return nil, fmt.Errorf("Config:>Failed to parse gcfg data: %s", err)
	}
	return &cfg, nil
}

func GetConfig() *Config {
	return &cfg
}

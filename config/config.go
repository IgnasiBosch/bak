package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".bak")
	configFile := filepath.Join(configDir, "config.json")

	// Return empty config if file doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Override with environment variables if present
	if env := os.Getenv("S3_ENDPOINT"); env != "" {
		config.Endpoint = env
	}
	if env := os.Getenv("S3_ACCESS_KEY"); env != "" {
		config.AccessKey = env
	}
	if env := os.Getenv("S3_SECRET_KEY"); env != "" {
		config.SecretKey = env
	}

	return &config, nil
}

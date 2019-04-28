package config

import "github.com/kelseyhightower/envconfig"

type DatabaseConfig struct {
	Path string `envconfig:"path" default:"/tmp/cursor.db"` // file path for database file
}

type Config struct {
	ServerURL string `envconfig:"server_url" default:"localhost:8080"`
	Debug     bool   `envconfig:"debug" default:"false"`
	Database  *DatabaseConfig
}

func FromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

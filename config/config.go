package config

import "github.com/kelseyhightower/envconfig"

// DatabaseConfig represents configurations for database assets of
type DatabaseConfig struct {
	Path string `envconfig:"cursor_db_path" default:"/tmp/cursor.db"` // file path for database file
}

// Config represents overall configuration objects of the application
type Config struct {
	ServerURL string `envconfig:"cursor_server_url" default:"localhost:8080"`
	Debug     bool   `envconfig:"cursor_debug" default:"false"`
	Database  *DatabaseConfig
}

// FromEnv parses environment variables into the config object based on envconfig name
func FromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

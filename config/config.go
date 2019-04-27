package config

type DatabaseConfig struct {
}

type Config struct {
	ServerURL string `envconfig:"server_url" default:"localhost:8080"`
	Debug     bool   `envconfig:"debug" default:"false"`
	Database  *DatabaseConfig
}

package config

import "github.com/kelseyhightower/envconfig"

// DatabaseConfig defines config settings for cusor database
type DatabaseConfig struct {
	// The database engine used by the backend
	// possible values are: googleDatastore
	Engine string `envconfig:"database_engine" default:"googleDatastore"`
	BoltDB *BoltDBConfig
}

// BoltDBConfig represents google firebase datastore configuration
// https://cloud.google.com/datastore/docs/concepts/overview
type BoltDBConfig struct {
	Path string `envconfig:"cursor_db_path" default:"/tmp/cursor.db"` // file path for database file
}

// MasterConfig defines config settings for the cursor master
type MasterConfig struct {
	IDLength            int    `envconfig:"cursor_master_id_length" default:"5"` // the length of all randomly generated ids
	HTTPURL             string `envconfig:"cursor_master_http_url" default:"localhost:8080"`
	GRPCURL             string `envconfig:"cursor_master_grpc_url" default:"localhost:8081"`
	RepoDirectoryPath   string `envconfig:"cursor_master_repo_directory_path" default:"/tmp/cursortest/repositories"`
	PluginDirectoryPath string `envconfig:"cursor_plugin_directory_path" default:"/tmp/cursortest/plugins"`
}

// FrontendConfig represents configuration for frontend basecoat
type FrontendConfig struct {
	Enable bool `envconfig:"frontend_enable" default:"true"`
	// This envvar is not used from this config but is here for completeness
	// it is set in the makefile and is injected into the js code at build time.
	// It controls where the frontend client should look for the gprc backend
	APIHost string `envconfig:"frontend_api_host" default:"https://localhost:8080"`
}

// CommandLineConfig represents configuration for cli application
type CommandLineConfig struct {
	Token string `envconfig:"token" default:""`
}

// Config represents overall configuration objects of the application
type Config struct {
	Debug       bool   `envconfig:"cursor_debug" default:"false"`
	TLSCertPath string `envconfig:"tls_cert_path" default:"./localhost.crt"`
	TLSKeyPath  string `envconfig:"tls_key_path" default:"./localhost.key"`
	Database    *DatabaseConfig
	Master      *MasterConfig
	CommandLine *CommandLineConfig
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

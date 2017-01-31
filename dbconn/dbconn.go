// Package dbconn creates database connection strings
// for "database/sql.Open(driverName string, dataSourceName string)" where
// dataSourceName is the connection string
package dbconn

import (
	"fmt"
	"os"
	"strings"
)

// Config has all fields needed to create a database connection
type Config struct {
	Database string
	User     string
	Password string
	Host     string
	Port     string
}

// FromENV creates a config from environment variables
// the prefix will be Joined with underscores to the env variables
func FromENV(prefixs ...string) Config {
	prefix := strings.Join(prefixs, "_")
	if prefix != "" {
		prefix = prefix + "_"
	}
	return Config{
		Database: os.Getenv(prefix + "DATABASE"),
		User:     os.Getenv(prefix + "USER"),
		Password: os.Getenv(prefix + "PASSWORD"),
		Host:     os.Getenv(prefix + "HOST"),
		Port:     os.Getenv(prefix + "PORT"),
	}
}

// MySQL returns a mysql conn string from Config's fields.
// Used for database/sql.Open("mysql", Config.MySQL())
func (cfg Config) MySQL() (conn string) {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s",
		cfg.Database, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

// TODO: More ouputs
// sqlite, postgres

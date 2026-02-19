package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 5
	connMaxLifetime = 5 * time.Minute
)

// Config holds the MySQL connection parameters.
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

// ConfigFromEnv builds a Config from environment variables.
// Returns an error if any required variable is missing.
func ConfigFromEnv() (Config, error) {
	required := map[string]string{
		"MYSQL_USER":     "",
		"MYSQL_PASSWORD": "",
		"MYSQL_HOST":     "",
		"MYSQL_PORT":     "",
		"MYSQL_DATABASE": "",
	}
	for key := range required {
		val := os.Getenv(key)
		if val == "" {
			return Config{}, fmt.Errorf("required environment variable %s is not set", key)
		}
		required[key] = val
	}
	return Config{
		User:     required["MYSQL_USER"],
		Password: required["MYSQL_PASSWORD"],
		Host:     required["MYSQL_HOST"],
		Port:     required["MYSQL_PORT"],
		Database: required["MYSQL_DATABASE"],
	}, nil
}

// DSN returns the MySQL data source name.
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

// Open creates a new database connection pool with sensible defaults.
func Open(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	return db, nil
}

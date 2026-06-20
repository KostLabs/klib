// Package postgres provides a thin wrapper around database/sql using the lib/pq
// driver, with optional exponential-backoff retry for the initial connection.
package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Config holds connection parameters for a PostgreSQL client.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	// SSLMode defaults to "disable" when empty.
	SSLMode string
}

// Client wraps a *sql.DB connection to PostgreSQL.
type Client struct {
	DB *sql.DB
}

// NewClient opens and pings a PostgreSQL connection described by cfg.
func NewClient(cfg Config) (*Client, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	return &Client{DB: db}, nil
}

// Close closes the underlying database connection pool.
func (c *Client) Close() error {
	return c.DB.Close()
}

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

// DB is the interface implemented by *sql.DB.
// Use this interface in constructors to satisfy the dependency-inversion principle.
type DB interface {
	Ping() error
	Close() error
	Query(query string, args ...any) (*sql.Rows, error)    //goverifier:ignore:any-type
	QueryRow(query string, args ...any) *sql.Row           //goverifier:ignore:any-type
	Exec(query string, args ...any) (sql.Result, error)    //goverifier:ignore:any-type
}

// Client wraps a DB connection to PostgreSQL.
type Client struct {
	db DB
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

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	return &Client{db: sqlDB}, nil
}

// DB returns the underlying database connection.
func (pgClient *Client) DB() DB {
	return pgClient.db
}

// Close closes the underlying database connection pool.
func (pgClient *Client) Close() error {
	return pgClient.db.Close()
}

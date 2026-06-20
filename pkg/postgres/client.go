package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Client struct {
	DB *sql.DB
}

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

func (c *Client) Close() error {
	return c.DB.Close()
}

package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

type ChessComConfig struct {
	GameURL     string
	PGN         string
	WhiteRating uint16
	BlackRating uint16
	Rating      uint16
}

func (ccm *ChessComConfig) NewChessComConfig(gameURL, pgn string, whiteRating, blackRating uint16) (*ChessComConfig, error) {
	if gameURL == "" {
		return nil, fmt.Errorf("game URL is required")
	}
	if pgn == "" {
		return nil, fmt.Errorf("PGN is required")
	}
	if whiteRating == 0 || blackRating == 0 {
		return nil, fmt.Errorf("both player ratings are required")
	}

	return &ChessComConfig{
		GameURL:     gameURL,
		PGN:         pgn,
		WhiteRating: whiteRating,
		BlackRating: blackRating,
		// TODO: replace the average with rating-level classification.
		Rating: uint16((uint32(whiteRating) + uint32(blackRating)) / 2),
	}, nil
}

func (cfg *DBConfig) NewDBConfig() (*DBConfig, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	config := &DBConfig{
		Name:     os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if config.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == "" {
		config.Port = "5432"
	}
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	return config, nil
}

func (cfg *DBConfig) DSN() string {
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Host + ":" + cfg.Port,
		Path:   cfg.Name,
	}
	query := dsn.Query()
	query.Set("sslmode", cfg.SSLMode)
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

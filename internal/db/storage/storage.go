package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Game is the database representation of a Chess.com game.
type Game struct {
	ID            int64
	PGN           string
	WhiteUsername string
	BlackUsername string
	PlayedAt      time.Time
	WhiteRating   uint16
	WhiteAccuracy *float64
	BlackRating   uint16
	BlackAccuracy *float64
	PlayerColor   string
	Result        string
	FEN           string
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) (*Repository, error) {
	if pool == nil {
		return nil, errors.New("database pool is nil")
	}
	return &Repository{pool: pool}, nil
}

const gameColumns = `
	id, pgn, white_username, black_username, played_at,
    white_rating, white_accuracy, black_rating, black_accuracy, 
	player_color, result, fen`

// SaveGame inserts a game.
// On success it sets game.ID to the stored row's ID.
func (r *Repository) SaveGame(ctx context.Context, game *Game) error {
	if game == nil {
		return errors.New("game is nil")
	}
	if game.PGN == "" {
		return errors.New("PGN is required")
	}
	if game.PlayedAt.IsZero() {
		return errors.New("played_at is required")
	}
	if game.WhiteRating == 0 || game.BlackRating == 0 {
		return errors.New("both player ratings are required")
	}

	err := r.pool.QueryRow(ctx, `
		INSERT INTO games (pgn, white_username, black_username, played_at,
			white_rating, white_accuracy, black_rating, black_accuracy, player_color, result, fen, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`,
		game.PGN, game.WhiteUsername, game.BlackUsername, game.PlayedAt,
		game.WhiteRating, game.WhiteAccuracy, game.BlackRating, game.BlackAccuracy,
		game.PlayerColor, game.Result, game.FEN,
		(int32(game.WhiteRating)+int32(game.BlackRating))/2,
	).Scan(&game.ID)
	if err != nil {
		return fmt.Errorf("save game: %w", err)
	}
	return nil
}

// GetGameByID returns pgx.ErrNoRows when no game has the requested ID.
func (r *Repository) GetGameByID(ctx context.Context, id int64) (*Game, error) {
	game, err := scanGame(r.pool.QueryRow(ctx,
		"SELECT "+gameColumns+" FROM games WHERE id = $1", id))
	if err != nil {
		return nil, fmt.Errorf("get game by ID: %w", err)
	}
	return game, nil
}

// ListGames sorts by a known column. Empty options default to played_at DESC.
func (r *Repository) ListGames(ctx context.Context, sortBy, sortOrder string) ([]Game, error) {
	allowed := map[string]string{
		"played_at": "played_at", "white_rating": "white_rating",
		"black_rating": "black_rating", "white_accuracy": "white_accuracy",
		"black_accuracy": "black_accuracy", "white_username": "white_username",
		"black_username": "black_username",
	}
	if sortBy == "" {
		sortBy = "played_at"
	}
	column, ok := allowed[sortBy]
	if !ok {
		return nil, fmt.Errorf("unsupported sort_by: %q", sortBy)
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	order := strings.ToUpper(sortOrder)
	if order != "ASC" && order != "DESC" {
		return nil, fmt.Errorf("unsupported sort_order: %q", sortOrder)
	}

	rows, err := r.pool.Query(ctx, "SELECT "+gameColumns+" FROM games ORDER BY "+column+" "+order+" NULLS LAST, id DESC")
	if err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	defer rows.Close()

	games := make([]Game, 0)
	for rows.Next() {
		game, err := scanGame(rows)
		if err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}
		games = append(games, *game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list games: %w", err)
	}
	return games, nil
}

type gameScanner interface {
	Scan(dest ...any) error
}

func scanGame(row gameScanner) (*Game, error) {
	game := new(Game)
	var whiteRating, blackRating int32
	if err := row.Scan(&game.ID, &game.PGN, &game.WhiteUsername,
		&game.BlackUsername, &game.PlayedAt, &whiteRating, &game.WhiteAccuracy,
		&blackRating, &game.BlackAccuracy, &game.PlayerColor, &game.Result, &game.FEN); err != nil {
		return nil, err
	}
	game.WhiteRating = uint16(whiteRating)
	game.BlackRating = uint16(blackRating)
	return game, nil
}

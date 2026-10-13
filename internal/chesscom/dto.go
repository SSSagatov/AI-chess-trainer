package chesscom

import "time"

// GameByURLRequest is the JSON body for requesting a Chess.com game by URL.
type GameByURLRequest struct {
	URL string `json:"url"`
}

// GameResponse contains the Game fields
type GameResponse struct {
	ID            int64     `json:"id"`
	PGN           string    `json:"pgn"`
	WhiteUsername string    `json:"white_username"`
	BlackUsername string    `json:"black_username"`
	PlayedAt      time.Time `json:"played_at"`
	WhiteRating   uint16    `json:"white_rating"`
	WhiteAccuracy *float64  `json:"white_accuracy"`
	BlackRating   uint16    `json:"black_rating"`
	BlackAccuracy *float64  `json:"black_accuracy"`
	PlayerColor   string    `json:"player_color"`
	Result        string    `json:"result"`
	FEN           string    `json:"fen"`
}

// GamesRequest defines JSON parameters for listing and sorting saved games.
// SortBy accepts "played_at", "white_rating", "black_rating",
// "white_accuracy", "black_accuracy", "white_username", or "black_username".
// SortOrder accepts "asc" or "desc".
type GamesRequest struct {
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// GamesResponse reports the ordering used for the returned games.
type GamesResponse struct {
	Games     []GameResponse `json:"games"`
	SortBy    string         `json:"sort_by"`
	SortOrder string         `json:"sort_order"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

package controller

// GameByURLRequest is the JSON body for requesting a Chess.com game by URL.
type GameByURLRequest struct {
	URL string `json:"url"`
}

// GameResponse contains the game data returned to the client.
type GameResponse struct {
	URL         string `json:"url"`
	PGN         string `json:"pgn"`
	WhiteRating uint16 `json:"white_rating"`
	BlackRating uint16 `json:"black_rating"`
	Rating      uint16 `json:"rating"`
}

// GamesRequest defines JSON parameters for listing and sorting saved games.
// SortBy accepts "rating", "white_rating", "black_rating", or "url".
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

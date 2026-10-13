package chesscom

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/SSSagatov/AI-chess-trainer/internal/db/storage"
	"github.com/jackc/pgx/v5"
)

type GameRepository interface {
	SaveGame(ctx context.Context, game *storage.Game) error
	GetGameByID(ctx context.Context, id int64) (*storage.Game, error)
	ListGames(ctx context.Context, sortBy, sortOrder string) ([]storage.Game, error)
}

// GameFetcher obtains and parses a Chess.com game
// Its implementation can be added separately.
type GameFetcher interface {
	FetchGame(ctx context.Context, gameURL string) (*storage.Game, error)
}

type Controller struct {
	repository GameRepository
	fetcher    GameFetcher
}

func NewController(repository GameRepository, fetcher GameFetcher) *Controller {
	return &Controller{repository: repository, fetcher: fetcher}
}

// CreateGameByURL accepts {"url":"..."}, fetches the game, and saves it.
func (c *Controller) CreateGameByURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if c.fetcher == nil {
		writeError(w, http.StatusNotImplemented, "Chess.com game fetcher is not configured")
		return
	}

	var request GameByURLRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request must contain one JSON object")
		return
	}
	if strings.TrimSpace(request.URL) == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	game, err := c.fetcher.FetchGame(r.Context(), request.URL)
	if err != nil || game == nil {
		writeError(w, http.StatusBadGateway, "could not fetch Chess.com game")
		return
	}
	if err := c.repository.SaveGame(r.Context(), game); err != nil {
		writeError(w, http.StatusInternalServerError, "could not save game")
		return
	}
	writeJSON(w, http.StatusCreated, gameResponse(game))
}

// GetGameByID reads a saved game. The route must provide the {id} path value.
func (c *Controller) GetGameByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid game id")
		return
	}
	game, err := c.repository.GetGameByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not get game")
		return
	}
	writeJSON(w, http.StatusOK, gameResponse(game))
}

// ListGames reads sort_by and sort_order from the query string.
func (c *Controller) ListGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	request := GamesRequest{
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	}
	if request.SortBy == "" {
		request.SortBy = "played_at"
	}
	if request.SortOrder == "" {
		request.SortOrder = "desc"
	}
	if !validSort(request.SortBy, request.SortOrder) {
		writeError(w, http.StatusBadRequest, "invalid sort_by or sort_order")
		return
	}
	games, err := c.repository.ListGames(r.Context(), request.SortBy, request.SortOrder)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list games")
		return
	}
	responses := make([]GameResponse, 0, len(games))
	for i := range games {
		responses = append(responses, gameResponse(&games[i]))
	}
	writeJSON(w, http.StatusOK, GamesResponse{
		Games: responses, SortBy: request.SortBy, SortOrder: request.SortOrder,
	})
}

func validSort(by, order string) bool {
	switch by {
	case "played_at", "white_rating", "black_rating", "white_accuracy",
		"black_accuracy", "white_username", "black_username":
	default:
		return false
	}
	return order == "asc" || order == "desc"
}

func gameResponse(game *storage.Game) GameResponse {
	return GameResponse{
		ID: game.ID, PGN: game.PGN,
		WhiteUsername: game.WhiteUsername, BlackUsername: game.BlackUsername,
		PlayedAt: game.PlayedAt, WhiteRating: game.WhiteRating,
		WhiteAccuracy: game.WhiteAccuracy, BlackRating: game.BlackRating,
		BlackAccuracy: game.BlackAccuracy, PlayerColor: game.PlayerColor,
		Result: game.Result, FEN: game.FEN,
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

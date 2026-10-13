package chesscom

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SSSagatov/AI-chess-trainer/internal/db/storage"
	"github.com/SSSagatov/AI-chess-trainer/templates"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	controller *Controller
}

func NewHandler(pool *pgxpool.Pool) (*Handler, error) {
	repository, err := storage.NewRepository(pool)
	if err != nil {
		return nil, fmt.Errorf("create game repository: %w", err)
	}
	return &Handler{controller: NewController(repository, nil)}, nil
}

func (h *Handler) PageHandler() (http.Handler, error) {
	pageHandler, err := templates.Handler(h.loadPage)
	if err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", NewRouter(h))
	mux.Handle("/", pageHandler)
	return mux, nil
}

func (h *Handler) loadPage(ctx context.Context, sortBy, sortOrder string) (templates.PageData, error) {
	if sortBy == "" {
		sortBy = "played_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	if !validSort(sortBy, sortOrder) {
		return templates.PageData{}, templates.ErrInvalidSort
	}

	games, err := h.controller.repository.ListGames(ctx, sortBy, sortOrder)
	if err != nil {
		return templates.PageData{}, fmt.Errorf("load games: %w", err)
	}
	data := templates.PageData{
		Games:     make([]templates.GameRow, 0, len(games)),
		GameCount: len(games), Rating: "—", Accuracy: "—",
		SortBy: sortBy, SortOrder: sortOrder,
	}
	var latest *storage.Game
	var accuracySum float64
	var accuracyCount int
	for i := range games {
		game := &games[i]
		data.Games = append(data.Games, templates.GameRow{
			ID: game.ID, WhiteUsername: game.WhiteUsername,
			BlackUsername: game.BlackUsername, PlayedAt: game.PlayedAt,
			WhiteRating: game.WhiteRating, BlackRating: game.BlackRating,
			Result: game.Result,
		})
		if latest == nil || game.PlayedAt.After(latest.PlayedAt) {
			latest = game
		}
		var accuracy *float64
		switch strings.ToLower(game.PlayerColor) {
		case "white":
			accuracy = game.WhiteAccuracy
		case "black":
			accuracy = game.BlackAccuracy
		}
		if accuracy != nil {
			accuracySum += *accuracy
			accuracyCount++
		}
	}
	if latest != nil {
		switch strings.ToLower(latest.PlayerColor) {
		case "white":
			data.Rating = fmt.Sprint(latest.WhiteRating)
		case "black":
			data.Rating = fmt.Sprint(latest.BlackRating)
		}
	}
	if accuracyCount > 0 {
		data.Accuracy = fmt.Sprintf("%.1f%%", accuracySum/float64(accuracyCount))
	}
	return data, nil
}

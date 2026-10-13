package chesscom

import (
	"fmt"
	"net/http"

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
	pageHandler, err := templates.Handler()
	if err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", NewRouter(h))
	mux.Handle("/", pageHandler)
	return mux, nil
}

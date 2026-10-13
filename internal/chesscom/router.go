package chesscom

import "net/http"

// NewRouter maps game endpoints to the controller handlers.
// POST /api/games accepts GameByURLRequest as JSON.
// GET /api/games accepts sort_by and sort_order query parameters.
// GET /api/games/{id} returns one saved game.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/games", h.controller.CreateGameByURL)
	mux.HandleFunc("GET /api/games", h.controller.ListGames)
	mux.HandleFunc("GET /api/games/{id}", h.controller.GetGameByID)
	return mux
}

package templates

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandlerRendersGamesAndCSS(t *testing.T) {
	handler, err := Handler(func(_ context.Context, by, order string) (PageData, error) {
		if by != "played_at" || order != "desc" {
			t.Fatalf("sort options = %q, %q", by, order)
		}
		return PageData{
			Games:     []GameRow{{ID: 7, WhiteUsername: "Alice", BlackUsername: "Bob", PlayedAt: time.Date(2026, 10, 14, 0, 0, 0, 0, time.UTC), WhiteRating: 1200, BlackRating: 1180, Result: "win"}},
			GameCount: 1, Rating: "1200", Accuracy: "83.5%", SortBy: by, SortOrder: order,
		}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	page := httptest.NewRecorder()
	handler.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/?sort_by=played_at&sort_order=desc", nil))
	if page.Code != http.StatusOK {
		t.Fatalf("page status = %d", page.Code)
	}
	for _, value := range []string{"Alice", "Bob", "14.10.2026", "/api/games/7", "83.5%"} {
		if !strings.Contains(page.Body.String(), value) {
			t.Errorf("page does not contain %q", value)
		}
	}

	css := httptest.NewRecorder()
	handler.ServeHTTP(css, httptest.NewRequest(http.MethodGet, "/assets/css/pages/dashboard.css", nil))
	if css.Code != http.StatusOK || !strings.Contains(css.Body.String(), ".chessboard") {
		t.Fatalf("CSS response: status %d", css.Code)
	}
}

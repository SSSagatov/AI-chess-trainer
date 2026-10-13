package templates

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

//go:embed index.html assets/css/*.css assets/css/pages/*.css
var files embed.FS

var ErrInvalidSort = errors.New("invalid game sort options")

type GameRow struct {
	ID            int64
	WhiteUsername string
	BlackUsername string
	PlayedAt      time.Time
	WhiteRating   uint16
	BlackRating   uint16
	Result        string
}

type PageData struct {
	Games     []GameRow
	GameCount int
	Rating    string
	Accuracy  string
	SortBy    string
	SortOrder string
}

type PageLoader func(context.Context, string, string) (PageData, error)

// LoadHTML parses the HTML page from the embedded templates directory.
func LoadHTML() (*template.Template, error) {
	return template.New("index.html").Funcs(template.FuncMap{
		"formatDate": func(value time.Time) string {
			return value.Format("02.01.2006")
		},
		"gameLabel": func(white, black string) string {
			return fmt.Sprintf("%s — %s", white, black)
		},
	}).ParseFS(files, "index.html")
}

// CSSHandler serves embedded CSS files under /assets/css/.
func CSSHandler() (http.Handler, error) {
	css, err := fs.Sub(files, "assets/css")
	if err != nil {
		return nil, err
	}
	return http.StripPrefix("/assets/css/", http.FileServer(http.FS(css))), nil
}

// Handler serves the dashboard at / and its CSS assets at /assets/css/.
func Handler(load PageLoader) (http.Handler, error) {
	if load == nil {
		return nil, fmt.Errorf("page loader is nil")
	}
	page, err := LoadHTML()
	if err != nil {
		return nil, err
	}
	css, err := CSSHandler()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/css/", css)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		data, err := load(r.Context(), r.URL.Query().Get("sort_by"), r.URL.Query().Get("sort_order"))
		if err != nil {
			if errors.Is(err, ErrInvalidSort) {
				http.Error(w, "invalid game sort options", http.StatusBadRequest)
				return
			}
			http.Error(w, "could not load games", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := page.Execute(w, data); err != nil {
			return
		}
	})
	return mux, nil
}

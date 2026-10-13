package templates

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed index.html assets/css/*.css assets/css/pages/*.css
var files embed.FS

// LoadHTML parses the HTML page from the embedded templates directory.
func LoadHTML() (*template.Template, error) {
	return template.ParseFS(files, "index.html")
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
func Handler() (http.Handler, error) {
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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := page.Execute(w, nil); err != nil {
			// The template is embedded and parsed at startup; execution does not use dynamic data.
			return
		}
	})
	return mux, nil
}

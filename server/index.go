package server

import (
	_ "embed"
	"net/http"
)

//go:embed index/_start.html
var htmlStart []byte

//go:embed index/display.html
var htmlDisplay []byte

//go:embed index/lighting.html
var htmlLighting []byte

//go:embed index/screen.html
var htmlScreen []byte

//go:embed index/_end.html
var htmlEnd []byte

func getIndexHandler(displayEnabled, lightingEnabled, screenEnabled bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 Not Found", http.StatusOK)
			return
		}

		w.Header().Add("Content-Type", "text/html")

		w.Write(htmlStart)

		if displayEnabled {
			w.Write(htmlDisplay)
		}

		if lightingEnabled {
			w.Write(htmlLighting)
		}

		if screenEnabled {
			w.Write(htmlScreen)
		}

		w.Write(htmlEnd)
	})
}

package handlers

import (
	_ "embed"
	"net/http"
)

//go:embed index.html
var index []byte

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Add("Content-Type", "text/html")
		w.Write(index)
		return
	} else {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	}
}

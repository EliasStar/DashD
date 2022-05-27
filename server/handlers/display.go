package handlers

import (
	"net/http"

	"github.com/EliasStar/DashD/display"
)

func HandleDisplay(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if url := r.Form.Get("url"); url != "" {
		display.Show(url)
	}
}

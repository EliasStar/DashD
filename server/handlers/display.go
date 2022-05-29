package handlers

import (
	"net/http"
	"strconv"

	"github.com/EliasStar/DashD/display"
)

func HandleDisplay(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if url := r.Form.Get("url"); url != "" {
		display.Show(url)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleResize(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	width, err := strconv.ParseUint(r.Form.Get("width"), 10, 32)
	if err != nil {
		return
	}

	height, err := strconv.ParseUint(r.Form.Get("height"), 10, 32)
	if err != nil {
		return
	}

	display.Resize(uint(width), uint(height))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

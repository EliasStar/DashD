package server

import (
	"net/http"
	"strconv"

	"github.com/EliasStar/DashD/display"
)

func handleDisplay(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if url := r.Form.Get("url"); url != "" {
		display.Show(url)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	posX, err := strconv.ParseUint(r.Form.Get("posX"), 10, 32)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	posY, err := strconv.ParseUint(r.Form.Get("posY"), 10, 32)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	display.Move(uint(posX), uint(posY))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleResize(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	width, err := strconv.ParseUint(r.Form.Get("width"), 10, 32)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	height, err := strconv.ParseUint(r.Form.Get("height"), 10, 32)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	display.Resize(uint(width), uint(height))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

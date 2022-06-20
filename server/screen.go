package server

import (
	"net/http"

	"github.com/EliasStar/DashD/screen"
)

func handlePower(w http.ResponseWriter, r *http.Request) {
	screen.PushPowerButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleSource(w http.ResponseWriter, r *http.Request) {
	screen.PushSourceButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleMenu(w http.ResponseWriter, r *http.Request) {
	screen.PushMenuButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handlePlus(w http.ResponseWriter, r *http.Request) {
	screen.PushPlusButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleMinus(w http.ResponseWriter, r *http.Request) {
	screen.PushMinusButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

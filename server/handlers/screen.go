package handlers

import (
	"net/http"

	"github.com/EliasStar/DashD/screen"
)

func HandlePower(w http.ResponseWriter, r *http.Request) {
	screen.PushPowerButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleSource(w http.ResponseWriter, r *http.Request) {
	screen.PushSourceButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleMenu(w http.ResponseWriter, r *http.Request) {
	screen.PushMenuButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandlePlus(w http.ResponseWriter, r *http.Request) {
	screen.PushPlusButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HandleMinus(w http.ResponseWriter, r *http.Request) {
	screen.PushMinusButton()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package handlers

import (
	"net/http"

	"github.com/EliasStar/DashD/screen"
)

func HandlePower(w http.ResponseWriter, r *http.Request) {
	screen.PushPowerButton()
}

func HandleSource(w http.ResponseWriter, r *http.Request) {
	screen.PushSourceButton()
}

func HandleMenu(w http.ResponseWriter, r *http.Request) {
	screen.PushMenuButton()
}

func HandlePlus(w http.ResponseWriter, r *http.Request) {
	screen.PushPlusButton()
}

func HandleMinus(w http.ResponseWriter, r *http.Request) {
	screen.PushMinusButton()
}

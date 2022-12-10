package server

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/EliasStar/DashD/lighting"
	. "github.com/EliasStar/DashD/log"
)

func handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	ErrorIf("HTTP::HandleConfig", json.NewEncoder(w).Encode(
		map[string][]map[string]uint{
			"channels": {
				{"channel": 1, "leds": lighting.Length()},
			},
		},
	))
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var dec io.Reader
	for _, v := range r.Form {

		if len(v) > 0 {
			dec = base64.NewDecoder(base64.StdEncoding, strings.NewReader(v[0]))
			break
		}

		return
	}

	ignore := make([]byte, 2)
	if n, err := dec.Read(ignore); err != nil || n < 2 {
		return
	}

	var colors []lighting.RGB
	rgb := make([]byte, 3)
	for {
		n, err := dec.Read(rgb)
		if err != nil || n < 3 {
			break
		}

		colors = append(colors, lighting.RGB{R: rgb[0], G: rgb[1], B: rgb[2]})
	}

	lighting.Render(colors)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	lighting.Render(make([]lighting.RGB, lighting.Length()))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

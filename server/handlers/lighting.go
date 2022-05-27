package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/EliasStar/DashD/lighting"
	. "github.com/EliasStar/DashD/log"
)

func HandleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	type channel struct {
		Index  uint `json:"channel"`
		Length uint `json:"leds"`
	}

	type channelList struct {
		List []channel `json:"channels"`
	}

	ErrorIf("HTTP::HandleConfig", json.NewEncoder(w).Encode(channelList{
		List: []channel{{1, lighting.Length()}},
	}))
}

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	dec := base64.NewDecoder(base64.StdEncoding, r.Body)

	ignore := make([]byte, 2)
	if n, err := dec.Read(ignore); err != nil || n != 2 {
		Error("HTTP::HandleUpdate", "Failed to read ignore bytes: n=", n, "err=", err)
		return
	}

	var colors []lighting.RGB
	rgb := make([]byte, 3)
	for {
		n, err := dec.Read(rgb)
		if err == io.EOF {
			break
		}

		if err != nil || n != 3 {
			Error("HTTP::HandleUpdate", "Failed to read color: n=", n, "err=", err)
			return
		}

		colors = append(colors, lighting.RGB{R: rgb[0], G: rgb[1], B: rgb[2]})
	}

	lighting.Render(colors)
}

func HandleReset(w http.ResponseWriter, r *http.Request) {
	lighting.Render(make([]lighting.RGB, lighting.Length()))
}

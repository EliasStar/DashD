package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/EliasStar/DashD/lighting"
	. "github.com/EliasStar/DashD/log"
	"github.com/EliasStar/DashD/server/handlers"
)

const tag = "Server"

var server http.Server
var handler http.ServeMux
var socket net.PacketConn

func init() {
	Info(tag, "Starting.")
	server.SetKeepAlivesEnabled(false)

	handler.HandleFunc("/", handlers.HandleIndex)

	handler.HandleFunc("/display", handlers.HandleDisplay)
	handler.HandleFunc("/resize", handlers.HandleResize)

	handler.HandleFunc("/config", handlers.HandleConfig)
	handler.HandleFunc("/update", handlers.HandleUpdate)
	handler.HandleFunc("/reset", handlers.HandleReset)

	handler.HandleFunc("/power", handlers.HandlePower)
	handler.HandleFunc("/source", handlers.HandleSource)
	handler.HandleFunc("/menu", handlers.HandleMenu)
	handler.HandleFunc("/plus", handlers.HandlePlus)
	handler.HandleFunc("/minus", handlers.HandleMinus)

	server.Handler = logRequests(&handler)
}

func logRequests(handler http.Handler) http.Handler {
	log := func(w http.ResponseWriter, r *http.Request) {
		Info(tag, r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(log)
}

func ListenHTTP(port uint) {
	Info(tag, "HTTP listening on:", port)

	server.Addr = ":" + strconv.FormatUint(uint64(port), 10)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		PanicIf(tag, err)
	}
}

func ListenUDP(port uint) {
	Info(tag, "UDP listening on:", port)

	var err error
	socket, err = net.ListenPacket("udp", ":"+strconv.FormatUint(uint64(port), 10))
	PanicIf(tag, err)

	var lastIndex uint8
	rgb := make([]byte, 2+lighting.Length()*3)
	for {
		n, _, err := socket.ReadFrom(rgb)
		if errors.Is(err, net.ErrClosed) {
			break
		}

		if n == 0 && err != nil {
			Error(tag, "Failed to read packet:", err)
			continue
		}

		if n < 5 {
			continue
		}

		if (rgb[0] > lastIndex) || ((lastIndex > 225) && (rgb[0] < 25)) {
			lastIndex = rgb[0]

			var colors []lighting.RGB

			for i := 2; i+2 < n; i += 3 {
				colors = append(colors, lighting.RGB{R: rgb[i], G: rgb[i+1], B: rgb[i+2]})
			}

			lighting.Render(colors)
		}
	}
}

func Destroy() {
	Info(tag, "Stopping.")
	ErrorIf(tag, server.Shutdown(context.Background()))
	if socket != nil {
		ErrorIf(tag, socket.Close())
	}
}

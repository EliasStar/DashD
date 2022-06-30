package server

import (
	"context"
	_ "embed"
	"net"
	"net/http"
	"strconv"
	"sync"

	. "github.com/EliasStar/DashD/log"
)

const tag = "Server"

var wg sync.WaitGroup
var server http.Server

func Init(port uint, displayEnabled, lightingEnabled, screenEnabled bool) {
	Info(tag, "Starting.")
	server.SetKeepAlivesEnabled(false)
	server.Addr = net.JoinHostPort("", strconv.FormatUint(uint64(port), 10))

	var handler http.ServeMux
	handler.Handle("/", getIndexHandler(displayEnabled, lightingEnabled, screenEnabled))

	if displayEnabled {
		handler.HandleFunc("/display", handleDisplay)
		handler.HandleFunc("/resize", handleResize)
	}

	if lightingEnabled {
		handler.HandleFunc("/config", handleConfig)
		handler.HandleFunc("/update", handleUpdate)
		handler.HandleFunc("/reset", handleReset)
	}

	if screenEnabled {
		handler.HandleFunc("/power", handlePower)
		handler.HandleFunc("/source", handleSource)
		handler.HandleFunc("/menu", handleMenu)
		handler.HandleFunc("/plus", handlePlus)
		handler.HandleFunc("/minus", handleMinus)
	}

	server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Info(tag, r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})

	wg.Add(1)
	go listen()
}

func listen() {
	Info(tag, "Listening on:", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		PanicIf(tag, err)
	}

	wg.Done()
}

func Destroy() {
	Info(tag, "Stopping.")

	ErrorIf(tag, server.Shutdown(context.Background()))

	wg.Wait()
}

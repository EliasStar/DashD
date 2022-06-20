package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/EliasStar/DashD/display"
	"github.com/EliasStar/DashD/lighting"
	"github.com/EliasStar/DashD/screen"
	"github.com/EliasStar/DashD/server"
	"github.com/EliasStar/DashD/socket"
)

var wg = new(sync.WaitGroup)
var signalChannel = make(chan os.Signal, 2)

func main() {
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	flag.Parse()

	if *displayEnabled {
		display.Init(*windowWidth, *windowHeight, *defaultUrl)
	}

	if *lightingEnabled {
		lighting.Init(*ledstripPin, *ledstripLength)

		wg.Add(1)
		go socket.Listen(*udpPort, wg)
	}

	if *screenEnabled {
		screen.Init(*powerPin, *sourcePin, *menuPin, *plusPin, *minusPin)
	}

	server.Init(*displayEnabled, *lightingEnabled, *screenEnabled)

	wg.Add(1)
	go server.Listen(*httpPort, wg)

	<-signalChannel

	server.Destroy()

	if *displayEnabled {
		display.Destroy()
	}

	if *lightingEnabled {
		lighting.Destroy()
		socket.Destroy()
	}

	if *screenEnabled {
		screen.Destroy()
	}

	wg.Wait()
}

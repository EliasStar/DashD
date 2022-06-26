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

var displayEnabled bool
var windowWidth uint
var windowHeight uint
var defaultUrl string

var lightingEnabled bool
var ledstripPin uint
var ledstripLength uint

var screenEnabled bool
var powerPin uint
var sourcePin uint
var menuPin uint
var plusPin uint
var minusPin uint

var httpPort uint
var udpPort uint

var wg = new(sync.WaitGroup)
var signalChannel = make(chan os.Signal, 2)

func init() {
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Usage = nil

	flag.BoolVar(&displayEnabled, "display_enabled", true, "Enable display module")
	flag.UintVar(&windowWidth, "window_width", 1024, "width of the window")
	flag.UintVar(&windowHeight, "window_height", 768, "height of the window")
	flag.StringVar(&defaultUrl, "default_url", "data:text/html;base64,PGgxPkRhc2hEPC9oMT4KPHA+TGlnaHR3ZWlnaHQgZGFlbW9uIGZvciBSYXNwYmVycnkgUGkgZHJpdmVuIGtpb3NrczwvcD4=", "initial website to load")

	flag.BoolVar(&lightingEnabled, "lighting_enabled", true, "Enable lighting module")
	flag.UintVar(&ledstripPin, "ledstrip_pin", 18, "pin used for data line of the led strip")
	flag.UintVar(&ledstripLength, "ledstrip_length", 62, "number of leds in the led strip")

	flag.BoolVar(&screenEnabled, "screen_enabled", true, "Enable screen module")
	flag.UintVar(&powerPin, "power_pin", 17, "pin connected to the power button of the screen")
	flag.UintVar(&sourcePin, "source_pin", 24, "pin connected to the source button of the screen")
	flag.UintVar(&menuPin, "menu_pin", 27, "pin connected to the menu button of the screen")
	flag.UintVar(&plusPin, "plus_pin", 22, "pin connected to the plus button of the screen")
	flag.UintVar(&minusPin, "minus_pin", 23, "pin connected to the minus button of the screen")

	flag.UintVar(&httpPort, "http_port", 80, "port used by the http server")
	flag.UintVar(&udpPort, "udp_port", 1872, "port used by the lighting socket")

	flag.Parse()
}

func main() {
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	if displayEnabled {
		display.Init(windowWidth, windowHeight, defaultUrl)
	}

	if lightingEnabled {
		lighting.Init(ledstripPin, ledstripLength)

		wg.Add(1)
		go socket.Listen(udpPort, wg)
	}

	if screenEnabled {
		screen.Init(powerPin, sourcePin, menuPin, plusPin, minusPin)
	}

	server.Init(displayEnabled, lightingEnabled, screenEnabled)

	wg.Add(1)
	go server.Listen(httpPort, wg)

	<-signalChannel

	server.Destroy()

	if displayEnabled {
		display.Destroy()
	}

	if lightingEnabled {
		lighting.Destroy()
		socket.Destroy()
	}

	if screenEnabled {
		screen.Destroy()
	}

	wg.Wait()
}

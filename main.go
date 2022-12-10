package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/EliasStar/DashD/display"
	"github.com/EliasStar/DashD/lighting"
	"github.com/EliasStar/DashD/screen"
	"github.com/EliasStar/DashD/server"
	"github.com/EliasStar/DashD/socket"
)

var displayEnabled bool
var browserPath, defaultUrl string
var windowPosX, windowPosY, windowWidth, windowHeight uint

var lightingEnabled bool
var ledstripPin, ledstripLength uint

var screenEnabled bool
var powerPin, sourcePin, menuPin, plusPin, minusPin uint

var serverEnabled bool
var serverPort uint

var socketEnabled bool
var socketPort uint

var signalChannel = make(chan os.Signal, 2)

func init() {
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Usage = nil

	flag.BoolVar(&displayEnabled, "display_enabled", true, "Enable display module")
	flag.StringVar(&browserPath, "browser_path", "/usr/bin/chromium-browser", "path to chromium-based browser executable")
	flag.StringVar(&defaultUrl, "default_url", "data:text/html;base64,PGgxPkRhc2hEPC9oMT4KPHA+TGlnaHR3ZWlnaHQgZGFlbW9uIGZvciBSYXNwYmVycnkgUGkgZHJpdmVuIGtpb3NrczwvcD4=", "initial website to load")
	flag.UintVar(&windowPosX, "window_x", 0, "x position of the window")
	flag.UintVar(&windowPosY, "window_y", 0, "y position of the window")
	flag.UintVar(&windowWidth, "window_width", 1920, "width of the window")
	flag.UintVar(&windowHeight, "window_height", 1080, "height of the window")

	flag.BoolVar(&lightingEnabled, "lighting_enabled", true, "Enable lighting module")
	flag.UintVar(&ledstripPin, "ledstrip_pin", 18, "pin used for data line of the led strip")
	flag.UintVar(&ledstripLength, "ledstrip_length", 100, "number of leds in the led strip")

	flag.BoolVar(&screenEnabled, "screen_enabled", true, "Enable screen module")
	flag.UintVar(&powerPin, "power_pin", 17, "pin connected to the power button of the screen")
	flag.UintVar(&sourcePin, "source_pin", 24, "pin connected to the source button of the screen")
	flag.UintVar(&menuPin, "menu_pin", 27, "pin connected to the menu button of the screen")
	flag.UintVar(&plusPin, "plus_pin", 22, "pin connected to the plus button of the screen")
	flag.UintVar(&minusPin, "minus_pin", 23, "pin connected to the minus button of the screen")

	flag.BoolVar(&serverEnabled, "server_enabled", true, "Enable server module")
	flag.UintVar(&serverPort, "server_port", 80, "port used by the http server")

	flag.BoolVar(&socketEnabled, "socket_enabled", true, "Enable socket module")
	flag.UintVar(&socketPort, "socket_port", 1872, "port used by the lighting socket")

	flag.Parse()
}

func main() {
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	if displayEnabled {
		display.Init(browserPath, defaultUrl, windowPosX, windowPosY, windowWidth, windowHeight)
		display.Notify(signalChannel)
		defer display.Destroy()
	}

	if lightingEnabled {
		lighting.Init(ledstripPin, ledstripLength)
		defer lighting.Destroy()

		if socketEnabled {
			socket.Init(socketPort)
			defer socket.Destroy()
		}
	}

	if screenEnabled {
		screen.Init(powerPin, sourcePin, menuPin, plusPin, minusPin)
		defer screen.Destroy()
	}

	if serverEnabled {
		server.Init(serverPort, displayEnabled, lightingEnabled, screenEnabled)
		defer server.Destroy()
	}

	<-signalChannel
}

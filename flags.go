package main

import (
	"flag"
	"os"
)

var displayEnabled = flag.Bool("display_enabled", true, "Enable display module")
var windowWidth = flag.Uint("window_width", 1024, "width of the window")
var windowHeight = flag.Uint("window_height", 768, "height of the window")
var defaultUrl = flag.String("default_url", "data:text/html;base64,PGgxPkRhc2hEPC9oMT4KPHA+TGlnaHR3ZWlnaHQgZGFlbW9uIGZvciBSYXNwYmVycnkgUGkgZHJpdmVuIGtpb3NrczwvcD4=", "initial website to load")

var lightingEnabled = flag.Bool("lighting_enabled", true, "Enable lighting module")
var ledstripPin = flag.Uint("ledstrip_pin", 18, "pin used for data line of the led strip")
var ledstripLength = flag.Uint("ledstrip_length", 62, "number of leds in the led strip")

var screenEnabled = flag.Bool("screen_enabled", true, "Enable screen module")
var powerPin = flag.Uint("power_pin", 17, "pin connected to the power button of the screen")
var sourcePin = flag.Uint("source_pin", 24, "pin connected to the source button of the screen")
var menuPin = flag.Uint("menu_pin", 27, "pin connected to the menu button of the screen")
var plusPin = flag.Uint("plus_pin", 22, "pin connected to the plus button of the screen")
var minusPin = flag.Uint("minus_pin", 23, "pin connected to the minus button of the screen")

var httpPort = flag.Uint("http_port", 80, "port used by the http server")
var udpPort = flag.Uint("udp_port", 1872, "port used by the lighting socket")

func init() {
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Usage = nil
}

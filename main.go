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
)

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.Usage = nil

	displayWidth := flag.Uint("display_width", 1024, "width of the display")
	displayHeight := flag.Uint("display_height", 768, "height of the display")
	displayUrl := flag.String("display_url", "data:text/html;base64,PGgxPkRhc2hEPC9oMT4KPHA+TGlnaHR3ZWlnaHQgZGFlbW9uIGZvciBSYXNwYmVycnkgUGkgZHJpdmVuIGtpb3NrczwvcD4=", "initial website to load")

	ledstripPin := flag.Uint("ledstrip_pin", 18, "pin used for data line of the led strip")
	ledstripLength := flag.Uint("ledstrip_length", 62, "number of leds in the led strip")

	powerPin := flag.Uint("power_pin", 17, "pin connected to the power button of the screen")
	sourcePin := flag.Uint("source_pin", 24, "pin connected to the source button of the screen")
	menuPin := flag.Uint("menu_pin", 27, "pin connected to the menu button of the screen")
	plusPin := flag.Uint("plus_pin", 22, "pin connected to the plus button of the screen")
	minusPin := flag.Uint("minus_pin", 23, "pin connected to the minus button of the screen")

	httpPort := flag.Uint("http_port", 80, "port used by the http server")
	udpPort := flag.Uint("udp_port", 1872, "port used by the lighting socket")

	flag.Parse()

	screen.Init(*powerPin, *sourcePin, *menuPin, *plusPin, *minusPin)
	lighting.Init(*ledstripPin, *ledstripLength)
	display.Init(*displayWidth, *displayHeight, *displayUrl)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		server.ListenHTTP(*httpPort)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		server.ListenUDP(*udpPort)
		wg.Done()
	}()

	<-signalChannel

	server.Destroy()
	display.Destroy()
	lighting.Destroy()
	screen.Destroy()

	wg.Wait()
}

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

	ledstripPin := flag.Uint("ledstrip_pin", 18, "width of the display")
	ledstripLength := flag.Uint("ledstrip_length", 62, "height of the display")

	powerPin := flag.Uint("power_pin", 17, "width of the display")
	sourcePin := flag.Uint("source_pin", 24, "height of the display")
	menuPin := flag.Uint("menu_pin", 27, "width of the display")
	plusPin := flag.Uint("plus_pin", 22, "height of the display")
	minusPin := flag.Uint("minus_pin", 23, "height of the display")

	httpPort := flag.Uint("http_port", 80, "port used by the http server")
	udpPort := flag.Uint("udp_port", 1872, "port used by the lighting socket")

	flag.Parse()

	screen.Init(*powerPin, *sourcePin, *menuPin, *plusPin, *minusPin)
	lighting.Init(*ledstripPin, *ledstripLength)
	display.Init(*displayWidth, *displayHeight)

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

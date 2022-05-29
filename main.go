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
)

var mainChan = make(chan any)
var sigChan = make(chan os.Signal, 3)

func main() {
	go onExitSignal()
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.Usage = nil

	httpPort := flag.Uint("p", 80, "port used by the http server")
	udpPort := flag.Uint("u", 1872, "udp port used by the lighting socket")

	displayWidth := flag.Uint("w", 1920, "width of the display")
	displayHeight := flag.Uint("h", 1080, "height of the display")

	flag.Parse()

	display.Resize(*displayWidth, *displayHeight)

	go server.ListenHTTP(*httpPort)
	go server.ListenUDP(*udpPort)

	<-mainChan
}

func onExitSignal() {
	<-sigChan

	server.Destroy()

	display.Destroy()
	lighting.Destroy()
	screen.Destroy()

	close(mainChan)
}

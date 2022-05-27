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

	httpPort := flag.Uint("http_port", 80, "used by the http server")
	udpPort := flag.Uint("udp_port", 1872, "used by the lighting socket")

	flag.Parse()

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

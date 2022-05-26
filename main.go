package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/EliasStar/DashD/display"
	"github.com/EliasStar/DashD/lighting"
	"github.com/EliasStar/DashD/screen"
)

var sigChan = make(chan os.Signal, 3)

func main() {
	go onExitSignal()
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer destroy()

	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.Usage = nil

	//httpPort := flag.Uint("http_port", 80, "used by the http server")
	//udpPort := flag.Uint("udp_port", 1872, "used by the lighting socket")

	flag.Parse()
}

func onExitSignal() {
	<-sigChan
	destroy()
	os.Exit(0)
}

func destroy() {
	lighting.Destroy()
	display.Destroy()
	screen.Destroy()
}

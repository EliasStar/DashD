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

var wg sync.WaitGroup

func main() {
	run()
	wg.Wait()
}

func run() {
	defer func() {
		server.Destroy()
		display.Destroy()
		lighting.Destroy()
		screen.Destroy()
	}()

	sigChan := make(chan os.Signal, 3)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.Usage = nil

	httpPort := flag.Uint("p", 80, "port used by the http server")
	udpPort := flag.Uint("u", 1872, "udp port used by the lighting socket")

	displayWidth := flag.Uint("w", 800, "width of the display")
	displayHeight := flag.Uint("h", 600, "height of the display")

	flag.Parse()

	go func() {
		wg.Add(1)
		display.Init(*displayWidth, *displayHeight)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		server.ListenHTTP(*httpPort)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		server.ListenUDP(*udpPort)
		wg.Done()
	}()

	<-sigChan
}

package main

import (
	"flag"
	"os"
)

func main() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.CommandLine.Init("DashD", flag.ExitOnError)
	flag.CommandLine.Usage = nil

	//httpPort := flag.Uint("http_port", 80, "used by the http server")
	//udpPort := flag.Uint("udp_port", 1872, "used by the lighting socket")

	flag.Parse()
}

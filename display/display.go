package display

import (
	"os"

	. "github.com/EliasStar/DashD/log"
)

const tag = "Display"

var chromium *Chromium

func Init(browserPath, url string, posX, posY, width, height uint) {
	Info(tag, "Starting.")

	var err error
	chromium, err = NewChromium(browserPath, url, posX, posY, width, height)
	PanicIf(tag, err)

	go func () {
		PanicIf(tag, chromium.StartConnectionHandler())
	}()

	PanicIf(tag, chromium.InitConnection())
}

func Show(url string) {
	Info(tag, "Now showing:", url)
	ErrorIf(tag, chromium.Load(url))
}

func Move(posX, posY uint) {
	Info(tag, "Changing window location to:", posX, "|", posX)
	ErrorIf(tag, chromium.SetPosition(posX, posY))
}

func Resize(width, height uint) {
	Info(tag, "Changing window size to:", width, "x", height)
	ErrorIf(tag, chromium.SetSize(width, height))
}

func Notify(channel chan<- os.Signal) {
	go func ()  {
		ErrorIf(tag, chromium.Wait())
		channel <- os.Interrupt
	}()
}

func Destroy() {
	Info(tag, "Stopping.")
	ErrorIf(tag, chromium.Kill())
}

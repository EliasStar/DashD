package display

import (
	"sync"

	. "github.com/EliasStar/DashD/log"
	"github.com/webview/webview"
)

const tag = "Display"

var wg sync.WaitGroup
var window webview.WebView

var currentWidth, currentHeight int
var currentUrl string

func Init(width, height uint, url string) {
	Info(tag, "Starting.")

	wg.Add(1)
	go run()
}

func run() {
	window = webview.New(false)
	window.SetTitle("DashD")
	window.SetSize(currentWidth, currentHeight, webview.Hint(webview.HintNone))
	window.Navigate(currentUrl)

	window.Dispatch(func() {
		window.Run()
		window.Destroy()
	})

	wg.Done()
}

func Show(url string) {
	Info(tag, "Now showing:", url)

	currentUrl = url

	window.Dispatch(func() {
		window.Navigate(url)
	})
}

func Resize(width, height uint) {
	Info(tag, "Changed window size to:", width, "x", height)

	currentWidth, currentHeight = int(width), int(height)

	window.Dispatch(func() {
		window.SetSize(currentWidth, currentHeight, webview.Hint(webview.HintNone))
	})
}

func Destroy() {
	Info(tag, "Stopping.")
	window.Terminate()
	wg.Done()
}

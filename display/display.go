package display

import (
	"sync"

	. "github.com/EliasStar/DashD/log"
	"github.com/webview/webview"
)

const tag = "Display"

var wg sync.WaitGroup
var stopChannel chan any

var window webview.WebView

var currentWidth, currentHeight int
var currentUrl string

func Init(width, height uint, url string) {
	Info(tag, "Starting.")

	stopChannel = make(chan any)

	currentWidth, currentHeight = int(width), int(height)
	currentUrl = url

	go run()
}

func run() {
	for {
		window = webview.New(false)

		wg.Add(1)
		window.Dispatch(func() {
			window.SetTitle("DashD")
			window.SetSize(currentWidth, currentHeight, webview.Hint(webview.HintNone))
			window.Navigate(currentUrl)

			window.Run()
			window.Destroy()

			wg.Done()
		})

		wg.Wait()

		select {
		case <-stopChannel:
			return

		default:
			Info(tag, "Restarting.")
		}
	}
}

func Show(url string) {
	Info(tag, "Now showing:", url)

	currentUrl = url

	window.Dispatch(func() {
		window.Navigate(currentUrl)
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

	close(stopChannel)
	window.Terminate()

	wg.Wait()
}

package display

import (
	. "github.com/EliasStar/DashD/log"
	"github.com/webview/webview"
)

const tag = "Display"

var window webview.WebView
var stopChannel = make(chan any)
var returnChannel = make(chan any)

func Init(width, height uint) {
	Info(tag, "Starting.")

	go func() {
		for {
			window = webview.New(false)
			window.SetTitle("DashD")
			window.SetSize(int(width), int(height), webview.Hint(webview.HintNone))

			window.Run()
			window.Destroy()

			select {
			case <-stopChannel:
				close(returnChannel)
				return

			default:
				Info(tag, "Restarting.")
			}
		}
	}()
}

func Show(url string) {
	Info(tag, "Now showing:", url)
	window.Dispatch(func() {
		window.Navigate(url)
	})
}

func Resize(width, height uint) {
	Info(tag, "Changed window size to:", width, "x", height)
	window.Dispatch(func() {
		window.SetSize(int(width), int(height), webview.Hint(webview.HintNone))
	})
}

func Destroy() {
	Info(tag, "Stopping.")
	close(stopChannel)
	window.Terminate()
	<-returnChannel
}

package display

import (
	"sync"

	. "github.com/EliasStar/DashD/log"
	"github.com/webview/webview"
)

const tag = "Display"

var window webview.WebView
var stopChannel = make(chan any)
var wg sync.WaitGroup

func init() {
	Info(tag, "Starting.")
	window = webview.New(false)

	wg.Add(1)
	go func() {
		for {
			window.Run()
			window.Destroy()

			select {
			case <-stopChannel:
				wg.Done()
				return

			default:
				Info(tag, "Restarting.")
				window = webview.New(false)
			}
		}
	}()
}

func Show(url string) {
	window.Dispatch(func() {
		window.Navigate(url)
		Info(tag, "Now showing:", url)
	})
}

func Resize(width, height uint) {
	window.Dispatch(func() {
		window.SetSize(int(width), int(height), webview.HintNone)
		Info(tag, "Changed window size to:", width, "x", height)
	})
}

func Destroy() {
	Info(tag, "Stopping.")
	close(stopChannel)
	window.Terminate()
	wg.Wait()
}

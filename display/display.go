package display

import (
	"sync"

	. "github.com/EliasStar/DashD/log"
	"github.com/webview/webview"
)

const tag = "DISPLAY"

var window webview.WebView
var stopChannel chan bool
var wg sync.WaitGroup

func init() {
	Info(tag, "Starting.")
	stopChannel = make(chan bool, 1)
	wg.Add(1)

	go func() {
		for {
			window = webview.New(false)
			window.SetTitle("DashD")

			window.Run()

			window.Destroy()

			select {
			case <-stopChannel:
				wg.Done()
				return

			default:
				Info(tag, "Restarting.")
			}
		}
	}()
}

// Displays the content from the given URL.
func Show(url string) {
	window.Navigate(url)
	Info(tag, "Now showing:", url)
}

// Stops the UI goroutine and disposes the webview.
func Destroy() {
	Info(tag, "Stopping.")
	stopChannel <- false
	window.Terminate()
	wg.Wait()
}

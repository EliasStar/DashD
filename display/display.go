package display

import (
	"sync"

	"github.com/webview/webview"
)

var window webview.WebView
var stopChannel chan bool
var wg sync.WaitGroup

func init() {
	stopChannel = make(chan bool)
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
			}
		}
	}()
}

// Displays the content from the given URL.
func Show(url string) {
	window.Navigate(url)
}

// Stops the UI goroutine and disposes the webview.
func Destroy() {
	stopChannel <- false
	wg.Wait()
}

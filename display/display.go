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

func Show(url string) {
	window.Navigate(url)
	Info(tag, "Now showing:", url)
}

func Destroy() {
	Info(tag, "Stopping.")
	close(stopChannel)
	window.Terminate()
	wg.Wait()
}

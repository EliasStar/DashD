package socket

import (
	"errors"
	"net"
	"strconv"
	"sync"

	"github.com/EliasStar/DashD/lighting"
	. "github.com/EliasStar/DashD/log"
)

const tag = "Socket"

var socket net.PacketConn

func Listen(port uint, wg *sync.WaitGroup) {
	Info(tag, "Listening on:", port)

	var err error
	socket, err = net.ListenPacket("udp", ":"+strconv.FormatUint(uint64(port), 10))
	PanicIf(tag, err)

	var lastIndex uint8
	rgb := make([]byte, 2+lighting.Length()*3)
	for {
		n, _, err := socket.ReadFrom(rgb)
		if errors.Is(err, net.ErrClosed) {
			break
		}

		if n == 0 && err != nil {
			Error(tag, "Failed to read packet:", err)
			continue
		}

		if n < 5 {
			continue
		}

		if (rgb[0] > lastIndex) || ((lastIndex > 225) && (rgb[0] < 25)) {
			lastIndex = rgb[0]

			var colors []lighting.RGB

			for i := 2; i+2 < n; i += 3 {
				colors = append(colors, lighting.RGB{R: rgb[i], G: rgb[i+1], B: rgb[i+2]})
			}

			lighting.Render(colors)
		}
	}

	wg.Done()
}

func Destroy() {
	Info(tag, "Stopping.")
	if socket != nil {
		ErrorIf(tag, socket.Close())
	}
}

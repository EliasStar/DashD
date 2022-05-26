package lighting

import (
	. "github.com/EliasStar/DashD/log"
	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

const tag = "LIGHTING"

var strip *ws281x.WS2811

func init() {
	Info(tag, "Starting.")
	opt := ws281x.DefaultOptions
	channel := &opt.Channels[0]
	channel.GpioPin = 18
	channel.LedCount = 63
	channel.Brightness = 255

	var err error
	strip, err = ws281x.MakeWS2811(&opt)
	PanicIf(tag, err)
}

func Render(colors []RGB) error {
	if err := strip.Wait(); err != nil {
		Error(tag, "Error while waiting for the last frame to render:", err)
		return err
	}

	leds := strip.Leds(0)[1:]

	length := len(colors)
	if length > len(leds) {
		length = len(leds)
	}

	for i := 0; i < length; i++ {
		leds[i] = colors[i].ToUint32()
	}

	return strip.Render()
}

func Length() int {
	return len(strip.Leds(0)) - 1
}

func Destroy() {
	Info(tag, "Stopping.")
	strip.Fini()
}

type RGB struct {
	R, G, B uint8
}

func (c RGB) ToUint32() (color uint32) {
	color |= 0xff000000
	color |= uint32(c.R) << 16
	color |= uint32(c.G) << 8
	color |= uint32(c.B)

	return
}

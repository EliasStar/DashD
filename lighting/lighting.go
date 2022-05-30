package lighting

import (
	. "github.com/EliasStar/DashD/log"
	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

const tag = "Lighting"

var strip *ws281x.WS2811

func Init(ledstripPin, ledstripLength uint) {
	Info(tag, "Starting.")

	opt := ws281x.DefaultOptions
	opt.Channels[0].GpioPin = int(ledstripPin)
	opt.Channels[0].LedCount = int(ledstripLength)
	opt.Channels[0].Brightness = 255

	var err error
	strip, err = ws281x.MakeWS2811(&opt)
	PanicIf(tag, err)

	PanicIf(tag, strip.Init())
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

func Length() uint {
	return uint(len(strip.Leds(0)) - 1)
}

func Destroy() {
	Info(tag, "Stopping.")
	ErrorIf(tag, Render(make([]RGB, Length())))
	ErrorIf(tag, strip.Wait())
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

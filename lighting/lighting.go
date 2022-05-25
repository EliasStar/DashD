package lighting

import (
	. "github.com/EliasStar/DashD/utils"
	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

var strip *ws281x.WS2811

func init() {
	opt := ws281x.DefaultOptions
	channel := &opt.Channels[0]
	channel.GpioPin = 18
	channel.LedCount = 63
	channel.Brightness = 255

	var err error
	strip, err = ws281x.MakeWS2811(&opt)
	PanicIf(err)
}

func Render(colors []uint32) error {
	if err := strip.SetLedsSync(0, colors); err != nil {
		return err
	}

	return strip.Render()
}

func Length() int {
	return len(strip.Leds(0)) - 1
}

func Destroy() {
	strip.Fini()
}

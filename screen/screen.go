package screen

import (
	"time"

	. "github.com/EliasStar/DashD/utils"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var (
	buttonPower  *gpiod.Line
	buttonMenu   *gpiod.Line
	buttonPlus   *gpiod.Line
	buttonMinus  *gpiod.Line
	buttonSource *gpiod.Line
)

func init() {
	chip, err := gpiod.NewChip("gpiochip0", gpiod.WithConsumer("DashD"), gpiod.AsOutput(), gpiod.WithBiasDisabled)
	defer PanicIf(chip.Close())
	PanicIf(err)

	buttonPower, err = chip.RequestLine(rpi.GPIO17)
	PanicIf(err)

	buttonMenu, err = chip.RequestLine(rpi.GPIO27)
	PanicIf(err)

	buttonPlus, err = chip.RequestLine(rpi.GPIO22)
	PanicIf(err)

	buttonMinus, err = chip.RequestLine(rpi.GPIO23)
	PanicIf(err)

	buttonSource, err = chip.RequestLine(rpi.GPIO24)
	PanicIf(err)
}

func PressPowerButton() error {
	return pressButton(buttonPower)
}

func PressSourceButton() error {
	return pressButton(buttonSource)
}

func PressMenuButton() error {
	return pressButton(buttonMenu)
}

func PressPlusButton() error {
	return pressButton(buttonPlus)
}

func PressMinusButton() error {
	return pressButton(buttonMinus)
}

func pressButton(btn *gpiod.Line) (err error) {
	if err = btn.SetValue(1); err != nil {
		btn.SetValue(0)
		return
	}

	time.Sleep(250 * time.Millisecond)

	if err = btn.SetValue(0); err != nil {
		btn.SetValue(0)
		return
	}

	time.Sleep(250 * time.Millisecond)
	return
}

func Destroy() {
	PanicIf(buttonPower.Close())
	PanicIf(buttonMenu.Close())
	PanicIf(buttonPlus.Close())
	PanicIf(buttonMinus.Close())
	PanicIf(buttonSource.Close())
}

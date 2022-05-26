package screen

import (
	"time"

	. "github.com/EliasStar/DashD/log"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

const tag = "SCREEN"

var (
	buttonPower  *gpiod.Line
	buttonMenu   *gpiod.Line
	buttonPlus   *gpiod.Line
	buttonMinus  *gpiod.Line
	buttonSource *gpiod.Line
)

func init() {
	Info(tag, "Starting.")
	chip, err := gpiod.NewChip("gpiochip0", gpiod.WithConsumer("DashD"), gpiod.AsOutput(), gpiod.WithBiasDisabled)
	defer PanicIf(tag, chip.Close())
	PanicIf(tag, err)

	buttonPower, err = chip.RequestLine(rpi.GPIO17)
	PanicIf(tag, err)

	buttonMenu, err = chip.RequestLine(rpi.GPIO27)
	PanicIf(tag, err)

	buttonPlus, err = chip.RequestLine(rpi.GPIO22)
	PanicIf(tag, err)

	buttonMinus, err = chip.RequestLine(rpi.GPIO23)
	PanicIf(tag, err)

	buttonSource, err = chip.RequestLine(rpi.GPIO24)
	PanicIf(tag, err)
}

func PushPowerButton() error {
	return pushButton(buttonPower)
}

func PushSourceButton() error {
	return pushButton(buttonSource)
}

func PushMenuButton() error {
	return pushButton(buttonMenu)
}

func PushPlusButton() error {
	return pushButton(buttonPlus)
}

func PushMinusButton() error {
	return pushButton(buttonMinus)
}

func pushButton(btn *gpiod.Line) (err error) {
	if err = btn.SetValue(1); err != nil {
		Error(tag, "Error while trying to press button:", err)
		btn.SetValue(0)
		return
	}

	time.Sleep(250 * time.Millisecond)

	if err = btn.SetValue(0); err != nil {
		Error(tag, "Error while trying to release button:", err)
		btn.SetValue(0)
		return
	}

	time.Sleep(250 * time.Millisecond)
	return
}

func Destroy() {
	Info(tag, "Stopping.")
	PanicIf(tag, buttonPower.Close())
	PanicIf(tag, buttonMenu.Close())
	PanicIf(tag, buttonPlus.Close())
	PanicIf(tag, buttonMinus.Close())
	PanicIf(tag, buttonSource.Close())
}

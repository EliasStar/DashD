package screen

import (
	"time"

	. "github.com/EliasStar/DashD/log"
	"github.com/warthog618/gpiod"
)

const tag = "Screen"

var (
	buttonPower  *gpiod.Line
	buttonSource *gpiod.Line
	buttonMenu   *gpiod.Line
	buttonPlus   *gpiod.Line
	buttonMinus  *gpiod.Line
)

func Init(powerPin, sourcePin, menuPin, plusPin, minusPin uint) {
	Info(tag, "Starting.")

	var err error
	options := []gpiod.LineReqOption{gpiod.WithConsumer("DashD"), gpiod.AsOutput(), gpiod.WithBiasDisabled}

	buttonPower, err = gpiod.RequestLine("gpiochip0", int(powerPin), options...)
	PanicIf(tag, err)

	buttonSource, err = gpiod.RequestLine("gpiochip0", int(sourcePin), options...)
	PanicIf(tag, err)

	buttonMenu, err = gpiod.RequestLine("gpiochip0", int(menuPin), options...)
	PanicIf(tag, err)

	buttonPlus, err = gpiod.RequestLine("gpiochip0", int(plusPin), options...)
	PanicIf(tag, err)

	buttonMinus, err = gpiod.RequestLine("gpiochip0", int(minusPin), options...)
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
	ErrorIf(tag, buttonPower.SetValue(0))
	ErrorIf(tag, buttonPower.Close())

	ErrorIf(tag, buttonSource.SetValue(0))
	ErrorIf(tag, buttonSource.Close())

	ErrorIf(tag, buttonMenu.SetValue(0))
	ErrorIf(tag, buttonMenu.Close())

	ErrorIf(tag, buttonPlus.SetValue(0))
	ErrorIf(tag, buttonPlus.Close())

	ErrorIf(tag, buttonMinus.SetValue(0))
	ErrorIf(tag, buttonMinus.Close())
}

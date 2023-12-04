//go:build avr

package driverutils

import (
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"machine"
	"tinygo.org/x/drivers/hd44780i2c"
)

func InitLCD() (lcd hd44780i2c.Device, err error) {
	err = machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
	})
	if err != nil {
		return
	}

	lcd = hd44780i2c.New(machine.I2C0, 0x27) // some modules have address 0x3F
	err = lcd.Configure(hd44780i2c.Config{
		Width:       16, // required
		Height:      2,  // required
		CursorOn:    false,
		CursorBlink: false,
	})

	return
}

func SerialAck() {
	utils.Must(machine.Serial.WriteByte(utils.AsciiACK))
}

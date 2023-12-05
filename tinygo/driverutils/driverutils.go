//go:build avr

package driverutils

import (
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"machine"
	"tinygo.org/x/drivers/hd44780i2c"
)

func InitLCD() (lcd hd44780i2c.Device) {
	err := machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
	})
	if err != nil {
		panic(err)
	}

	lcd = hd44780i2c.New(machine.I2C0, 0x27) // some modules have address 0x3F
	err = lcd.Configure(hd44780i2c.Config{
		Width:       16, // required
		Height:      2,  // required
		CursorOn:    false,
		CursorBlink: false,
	})
	if err != nil {
		panic(err)
	}

	return
}

func SerialByte(b byte) {
	utils.Must(machine.Serial.WriteByte(b))
}

type SerialComm struct {
	readFromCurrFrame int
}

func (sc *SerialComm) Read(buf []byte) int {
	n, err := machine.Serial.Read(buf)
	utils.Must(err)

	if sc.readFromCurrFrame+n < utils.SerialPacketSize {
		sc.readFromCurrFrame += n
	} else if sc.readFromCurrFrame+n == utils.SerialPacketSize {
		SerialByte(utils.AsciiACK)
		sc.readFromCurrFrame = 0
	} else {
		print(utils.AsciiUS, "tx error. in packet: ", sc.readFromCurrFrame+n, utils.AsciiUS)
		panic("tx error")
	}

	return n
}

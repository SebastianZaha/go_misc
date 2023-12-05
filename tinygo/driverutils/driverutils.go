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

type TwoIntsLcd struct {
	device hd44780i2c.Device
	intBuf [9]byte

	// don't print too often, it is very slow, the driver calls sleep
	onceEvery  uint32
	lastWrote  uint32
	printCalls uint32
}

func (lcd *TwoIntsLcd) Init(onceEvery uint32) {
	lcd.onceEvery = onceEvery
	lcd.device = InitLCD()
	lcd.device.SetCursor(0, 0)
	lcd.device.Print([]byte("res: "))
	lcd.device.SetCursor(0, 1)
	lcd.device.Print([]byte("i: "))
}

func (lcd *TwoIntsLcd) OK(finalValue uint32) {
	lcd.device.SetCursor(15, 0)
	lcd.device.Print([]byte{'o'})
	lcd.device.SetCursor(15, 1)
	lcd.device.Print([]byte{'k'})

	// the final print should not be counted as "iteration"
	lcd.printCalls--
	lcd.Print(finalValue, false)
}

func (lcd *TwoIntsLcd) Print(n uint32, throttled bool) {
	lcd.printCalls++
	if throttled && lcd.printCalls > 1 && lcd.printCalls < lcd.lastWrote+lcd.onceEvery {
		return
	}
	lcd.lastWrote = lcd.printCalls

	lcd.device.SetCursor(5, 0)
	utils.FormatUint32(n, lcd.intBuf[:])
	lcd.device.Print(lcd.intBuf[:])
	lcd.device.SetCursor(5, 1)
	utils.FormatUint32(lcd.printCalls, lcd.intBuf[:])
	lcd.device.Print(lcd.intBuf[:])
}

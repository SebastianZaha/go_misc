//go:build avr

package main

import (
	"github.com/SebastianZaha/go_misc/tinygo/driverutils"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"tinygo.org/x/drivers/hd44780i2c"
)

var (
	// serial comms
	serialComm = &driverutils.SerialComm{}
	serialBuf  = make([]byte, utils.SerialPacketSize)

	// lcd and lcd utils
	lcd    hd44780i2c.Device
	intBuf = make([]byte, 8)

	txtSum   = []byte("Sum: ")
	txtTerms = []byte("of:  ")
	txtDone  = []byte("ok")

	// problem
	sum       uint32
	sumTerms  uint32
	prevDigit uint32
	lastWrote uint32
)

func main() {
	lcd = driverutils.InitLCD()

	lcd.SetCursor(0, 0)
	lcd.Print(txtSum)
	lcd.SetCursor(0, 1)
	lcd.Print(txtTerms)

	driverutils.SerialByte(utils.AsciiACK) // ready for input

	for {
		n := serialComm.Read(serialBuf)
		for i := 0; i < n; i++ {
			if serialBuf[i] == utils.AsciiEOT {
				lcd.SetCursor(14, 1)
				lcd.Print(txtDone)
				driverutils.SerialByte(utils.AsciiEOT)
				return
			} else if serialBuf[i] >= utils.Ascii0 && serialBuf[i] <= utils.Ascii9 {
				if prevDigit == 0 {
					prevDigit = uint32(serialBuf[i])
					add(10 * (prevDigit - 48))
				} else {
					prevDigit = uint32(serialBuf[i])
				}
			} else if serialBuf[i] == utils.AsciiUS { // new line in test input
				if prevDigit != 0 {
					add(prevDigit - 48)
					prevDigit = 0
				} // else a line with no digits on it
			} else {
				// irelevant char
			}
		}
	}
}

func add(n uint32) {
	sum += n
	sumTerms++

	if sumTerms < lastWrote+100 {
		return
	}
	lastWrote = sumTerms

	lcd.SetCursor(5, 0)
	utils.FormatUint32(sum, intBuf)
	lcd.Print(intBuf)
	lcd.SetCursor(5, 1)
	utils.FormatUint32(sumTerms, intBuf)
	lcd.Print(intBuf)
}

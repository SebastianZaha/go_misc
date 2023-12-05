//go:build avr

package main

import (
	"github.com/SebastianZaha/go_misc/tinygo/driverutils"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
)

var (
	serialComm driverutils.SerialComm
	serialBuf  = make([]byte, utils.SerialPacketSize)
	lcd        driverutils.TwoIntsLcd

	sum       uint32
	prevDigit uint32
)

func main() {
	lcd.Init(100)
	driverutils.SerialByte(utils.AsciiACK) // ready for input

	for {
		n := serialComm.Read(serialBuf)
		for i := 0; i < n; i++ {
			if serialBuf[i] == utils.AsciiEOT {
				lcd.OK(sum)
				driverutils.SerialByte(utils.AsciiEOT)
				return
			} else if serialBuf[i] >= utils.Ascii0 && serialBuf[i] <= utils.Ascii9 {
				if prevDigit == 0 {
					prevDigit = uint32(serialBuf[i])
					sum += 10 * (prevDigit - 48)
					lcd.Print(sum, true)
				} else {
					prevDigit = uint32(serialBuf[i])
				}
			} else if serialBuf[i] == utils.AsciiUS { // new line in test input
				if prevDigit != 0 {
					sum += prevDigit - 48
					lcd.Print(sum, true)
					prevDigit = 0
				} // else a line with no digits on it
			} else {
				// irelevant char
			}
		}
	}
}

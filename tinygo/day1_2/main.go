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
	digit     uint32

	digits = [][]byte{
		{'o', 'n', 'e'},
		{'t', 'w', 'o'},
		{'t', 'h', 'r', 'e', 'e'},
		{'f', 'o', 'u', 'r'},
		{'f', 'i', 'v', 'e'},
		{'s', 'i', 'x'},
		{'s', 'e', 'v', 'e', 'n'},
		{'e', 'i', 'g', 'h', 't'},
		{'n', 'i', 'n', 'e'},
	}
	digitTrack = []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
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
				digit = uint32(serialBuf[i])
				if prevDigit == 0 {
					sum += 10 * (digit - 48)
					lcd.Print(sum, true)
				}
				prevDigit = digit
				for j := 0; j < len(digitTrack); j++ {
					digitTrack[j] = 0
				}
			} else if serialBuf[i] == utils.AsciiUS { // new line in test input
				if prevDigit != 0 {
					sum += prevDigit - 48
					lcd.Print(sum, true)
					prevDigit = 0
				} // else a line with no digits on it
				for j := 0; j < len(digitTrack); j++ {
					digitTrack[j] = 0
				}
			} else {
				// letter
				j := 0
				for j = 0; j < len(digitTrack); j++ {
					if serialBuf[i] == digits[j][digitTrack[j]] {
						if digitTrack[j]+1 == len(digits[j]) {
							digit := uint32(j + 1)
							if prevDigit == 0 {
								sum += 10 * digit
								lcd.Print(sum, true)
							}
							prevDigit = digit + 48
							digitTrack[j] = 0
						} else {
							digitTrack[j]++
						}
					} else {
						if serialBuf[i] == digits[j][0] {
							digitTrack[j] = 1
						} else {
							digitTrack[j] = 0
						}
					}
				}
			}
		}
	}
}

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
				lcdPrint()
				lcd.Print(txtDone)
				driverutils.SerialByte(utils.AsciiEOT)
				return
			} else if serialBuf[i] >= utils.Ascii0 && serialBuf[i] <= utils.Ascii9 {
				digit = uint32(serialBuf[i])
				if prevDigit == 0 {
					add(10 * (digit - 48))
				}
				prevDigit = digit
				for j := 0; j < len(digitTrack); j++ {
					digitTrack[j] = 0
				}
			} else if serialBuf[i] == utils.AsciiUS { // new line in test input
				if prevDigit != 0 {
					add(prevDigit - 48)
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
								add(10 * digit)
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

func add(n uint32) {
	sum += n
	sumTerms++

	/*if sumTerms%2 == 0 {
		print(n, " = ", sum)
		driverutils.SerialByte(utils.AsciiUS)
	} else {
		driverutils.SerialByte(utils.AsciiUS)
		print(n/10, " + ")
	}*/
	if sumTerms < lastWrote+100 {
		return
	}
	lastWrote = sumTerms

	lcdPrint()
}

func lcdPrint() {
	lcd.SetCursor(5, 0)
	utils.FormatUint32(sum, intBuf)
	lcd.Print(intBuf)
	lcd.SetCursor(5, 1)
	utils.FormatUint32(sumTerms, intBuf)
	lcd.Print(intBuf)
}

//go:build avr

package main

import (
	"github.com/SebastianZaha/go_misc/tinygo/driverutils"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"tinygo.org/x/drivers/hd44780i2c"
)

type parseState byte

const (
	game parseState = iota
	colorNum
	colorName
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
	lastWrote uint32

	parsing                parseState
	currGame, currColorNum uint32

	maxR, maxG, maxB uint32
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
			dbg("#", serialBuf[i])
			if serialBuf[i] == utils.AsciiEOT {
				// finished parsing the input

				lcdPrint()
				lcd.SetCursor(14, 1)
				lcd.Print(txtDone)
				driverutils.SerialByte(utils.AsciiEOT)
				return
			} else if serialBuf[i] == utils.AsciiUS {
				// new line in test input

				add(maxR * maxG * maxB)

				parsing = game
				currGame = 0
				currColorNum = 0
				maxR = 0
				maxG = 0
				maxB = 0
			} else if serialBuf[i] >= utils.Ascii0 && serialBuf[i] <= utils.Ascii9 {
				if parsing == game {
					currGame = currGame*10 + uint32(serialBuf[i]-48)
				} else if parsing == colorNum {
					currColorNum = currColorNum*10 + uint32(serialBuf[i]-48)
				} else {
					dbg("got a digit. illegal parse state", serialBuf[i])
					return
				}
			} else if serialBuf[i] == ' ' {
				// If we expect to parse a number:
				//   - if we already parsed some digits, the space means the parsing is done.
				//   - if we didn't already parse digits, it could be some random space before the number
				if parsing == colorNum && currColorNum > 0 {
					parsing = colorName
				}
			} else if serialBuf[i] == ';' {
			} else if serialBuf[i] == ':' {
				parsing = colorNum
			} else if serialBuf[i] == ',' {
			} else if parsing == colorName {
				if serialBuf[i] == 'r' {
					if currColorNum > maxR {
						maxR = currColorNum
					}
					parsing = colorNum
					currColorNum = 0

				} else if serialBuf[i] == 'g' {
					if currColorNum > maxG {
						maxG = currColorNum
					}
					parsing = colorNum
					currColorNum = 0

				} else if serialBuf[i] == 'b' {
					if currColorNum > maxB {
						maxB = currColorNum
					}
					parsing = colorNum
					currColorNum = 0

				}
			}
		}
	}
}

func dbg(txt string, c byte) {
	return
	driverutils.SerialByte(utils.AsciiUS)
	print(txt, " '", c, "' p ", parsing, " s ", sum, " sT ", sumTerms, " cG ", currGame, " cC ", currColorNum)
	driverutils.SerialByte(utils.AsciiUS)
}

func add(n uint32) {
	sum += n
	sumTerms++
	if sumTerms < lastWrote+10 {
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

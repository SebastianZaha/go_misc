//go:build avr

package main

import (
	"github.com/SebastianZaha/go_misc/tinygo/driverutils"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
)

type parseState byte

const (
	game parseState = iota
	colorNum
	colorName
	impossible
)

const (
	maxR byte = 12
	maxG      = 13
	maxB      = 14
)

var (
	serialComm driverutils.SerialComm
	serialBuf  = make([]byte, utils.SerialPacketSize)
	lcd        driverutils.TwoIntsLcd

	sum                    uint32
	parsing                parseState
	currGame, currColorNum byte
)

func main() {
	lcd.Init(10)
	driverutils.SerialByte(utils.AsciiACK) // ready for input

	for {
		n := serialComm.Read(serialBuf)
		for i := 0; i < n; i++ {
			dbg("#", serialBuf[i])
			if serialBuf[i] == utils.AsciiEOT {
				// finished parsing the input
				lcd.OK(sum)
				driverutils.SerialByte(utils.AsciiEOT)
				return
			} else if serialBuf[i] == utils.AsciiUS {
				// new line in test input

				if parsing != impossible {
					sum += uint32(currGame)
					lcd.Print(sum, true)
				}

				parsing = game
				currGame = 0
				currColorNum = 0
			} else if parsing == impossible {
				continue
			} else if serialBuf[i] >= utils.Ascii0 && serialBuf[i] <= utils.Ascii9 {
				if parsing == game {
					currGame = currGame*10 + (serialBuf[i] - 48)
				} else if parsing == colorNum {
					currColorNum = currColorNum*10 + (serialBuf[i] - 48)
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
			} else if serialBuf[i] == ':' {
				parsing = colorNum
			} else if parsing == colorName {
				if serialBuf[i] == 'r' {
					if currColorNum > maxR {
						parsing = impossible
					} else {
						parsing = colorNum
						currColorNum = 0
					}
				} else if serialBuf[i] == 'g' {
					if currColorNum > maxG {
						parsing = impossible
					} else {
						parsing = colorNum
						currColorNum = 0
					}
				} else if serialBuf[i] == 'b' {
					if currColorNum > maxB {
						parsing = impossible
					} else {
						parsing = colorNum
						currColorNum = 0
					}
				}
			}
		}
	}
}

func dbg(txt string, c byte) {
	return
	driverutils.SerialByte(utils.AsciiUS)
	print(txt, " '", c, "' p ", parsing, " s ", sum, " cG ", currGame, " cC ", currColorNum)
	driverutils.SerialByte(utils.AsciiUS)
}

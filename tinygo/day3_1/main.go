//go:build avr

package main

import (
	"github.com/SebastianZaha/go_misc/tinygo/driverutils"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
)

type num struct {
	off    byte // offset in the row
	length byte // number of digits
}

var (
	lcd        driverutils.TwoIntsLcd
	serialComm driverutils.SerialComm
	serialBuf  = make([]byte, utils.SerialPacketSize)

	// Since we cannot process a full matrix, and we don't know the row length at compile time,
	// we will read a number of maxRows rows in a linear buffer.
	// Row 'i' will have offset rowLen*i.
	// Since we do not want to shift all the bytes around, when reaching maxRows, we have to
	// wrap around and overwrite previous rows.
	//
	// The api for handling these would depend on the necessary access patterns.
	//
	// Specifically for this problem, we parse rows as we read them,
	// we just need prevRow as we go along.
	// Numbers that *seem* to not have a symbol around *might* have one on the next row,
	// so their positions are saved.
	//
	// Once a row is completely read, *before* starting to read, we look at numbers on prevRow,
	// and check against currentRow (that we just finished reading).
	// When that is done, prevRow can be overwritten, currRow becomes prevRow,
	// and the numbers saved from prevRow are also nulled.
	//
	// Hence, we just need to buffer 2 rows at one time.

	data [800]byte

	rowLen uint16

	row0, row1, row2 []byte

	// index of the currently parsing row in the total number of rows in the input
	rowsParsed, curRowIdx uint16

	sum uint32

	curParsedNum     uint32
	curParsedNumGood bool
)

func special(row []byte, i uint16) bool {
	return row != nil && row[i] != '.' && (row[i] < utils.Ascii0 || row[i] > utils.Ascii9)
}

// We read the input in 3 row 'windows', and always process the middle of the 3.
func processRow1() {
	curParsedNum = 0
	curParsedNumGood = false
	var i uint16
	for i = 0; i < rowLen; i++ {
		if row1[i] >= utils.Ascii0 && row1[i] <= utils.Ascii9 {
			curParsedNum = curParsedNum*10 + uint32(row1[i]-utils.Ascii0)
			if curParsedNumGood {
				continue
			}
			if i > 0 {
				// previous column, on all 3 rows
				curParsedNumGood =
					special(row0, i-1) ||
						// this is only relevant for first digits of numbers, is false afterwards
						special(row1, i-1) ||
						special(row2, i-1)
			}
		} else if row1[i] == '.' {
			if curParsedNum > 0 {
				dbg2("fn", curParsedNum)
				if !curParsedNumGood {
					// prev and current column, on the other 2 rows than self (.)
					curParsedNumGood = special(row0, i-1) || special(row2, i-1) ||
						special(row0, i) || special(row2, i)
				}
				if curParsedNumGood {
					dbg2("a", curParsedNum)
					sum += curParsedNum
					lcd.Print(sum, true)
				}
				curParsedNum = 0
				curParsedNumGood = false
			}
		} else {
			// special symbol
			dbg("s", row1[i])

			if curParsedNum > 0 {
				dbg2("add", curParsedNum)
				sum += curParsedNum
				lcd.Print(sum, true)
				curParsedNum = 0
				curParsedNumGood = false
			}
		}
	}
	// did a number finish on last column?
	if curParsedNum > 0 {
		dbg2("frn", curParsedNum)

		if !curParsedNumGood {
			// last column, on the other 2 rows than self
			curParsedNumGood = special(row0, rowLen-1) || special(row2, rowLen)
		}
		if curParsedNumGood {
			dbg2("add", curParsedNum)
			sum += curParsedNum
			lcd.Print(sum, true)
		}
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			driverutils.SerialByte(utils.AsciiUS)
			println("panic ", err)
			driverutils.SerialByte(utils.AsciiUS)
		}
	}()
	driverutils.SerialByte(utils.AsciiUS)
	println("test")
	driverutils.SerialByte(utils.AsciiUS)

	row0[1000] = 2 // panic

	lcd.Init(10)
	driverutils.SerialByte(utils.AsciiACK) // ready for input

	row0 = nil
	row1 = nil
	row2 = data[:]

	for {
		n := serialComm.Read(serialBuf)
		for i := 0; i < n; i++ {
			if serialBuf[i] == utils.AsciiEOT {
				dbg("fin", serialBuf[i])
				// finished parsing the input
				row2 = nil
				processRow1()

				driverutils.SerialByte(utils.AsciiEOT)
				lcd.OK(sum)
				return
			} else if serialBuf[i] == utils.AsciiUS {
				dbg("eoln", serialBuf[i])
				rowsParsed++

				// new line in test input
				// is this the first line in the data that we finished reading?
				if row1 == nil {
					dbg("f0", 0)
					rowLen = curRowIdx
					row1 = data[:rowLen]
					row2 = data[rowLen : rowLen*2]
				} else if row0 == nil {
					dbg("f1", 0)
					processRow1()
					row0 = row1
					row1 = row2
					row2 = data[rowLen*2 : rowLen*3]
				} else {
					dbg("f2+", 0)
					processRow1()
					tmp := row0
					row0 = row1
					row1 = row2
					row2 = tmp
				}
				curRowIdx = 0
			} else {
				dbg("char", serialBuf[i])
				row2[curRowIdx] = serialBuf[i]
				curRowIdx++
			}
		}
	}
}

func dbg(s string, c byte) {
	driverutils.SerialByte(utils.AsciiUS)
	print(s, " '", c, "' cri ", curRowIdx, " row ", rowsParsed, " s ", sum, " rl ", rowLen)
	driverutils.SerialByte(utils.AsciiUS)
}

func dbg2(s string, i uint32) {
	driverutils.SerialByte(utils.AsciiUS)
	print(s, " ", i)
	driverutils.SerialByte(utils.AsciiUS)
}

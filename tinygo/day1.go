package main

import (
	"github.com/SebastianZaha/go_misc/tinygo/driverutils"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"machine"
	"tinygo.org/x/drivers/hd44780i2c"
)

var (
	sum      uint32
	sumTerms uint32
	lcd      hd44780i2c.Device

	buf    = make([]byte, 128) // utils.SerialPacketSize
	intBuf = make([]byte, 8)

	txtSum   = []byte("Sum: ")
	txtTerms = []byte("of:  ")
	txtDone  = []byte("ok")
)

func main() {
	var err error

	lcd, err = driverutils.InitLCD()
	utils.Must(err)

	lcd.SetCursor(0, 0)
	lcd.Print(txtSum)
	lcd.SetCursor(0, 1)
	lcd.Print(txtTerms)

	driverutils.SerialAck()

	var prevDigit uint32
	var readFromCurrFrame int
	for {
		n, err := machine.Serial.Read(buf)
		utils.Must(err)
		if n == 0 {
			continue
		}
		/*
			driverutils.SerialPacket()
			print("read bytes: ", n)
			driverutils.SerialPacket()
			print(string(buf[:n]))
			driverutils.SerialPacket()
		*/
		for i := 0; i < n; i++ {
			if buf[i] == utils.AsciiEOT {
				lcd.SetCursor(14, 1)
				lcd.Print(txtDone)
				driverutils.SerialPacket()
				print("done")
				driverutils.SerialPacket()
				return
			} else if buf[i] >= utils.Ascii0 && buf[i] <= utils.Ascii9 {
				if prevDigit == 0 {
					prevDigit = uint32(buf[i])
					add(10 * (prevDigit - 48))
				} else {
					prevDigit = uint32(buf[i])
				}
				//utils.SerialDbg("digit: ", prevDigit-48, "; curr sum: ", sum)
			} else if buf[i] == utils.AsciiUS {
				if prevDigit != 0 {
					add(prevDigit - 48)
					prevDigit = 0
				} // else a newline without a line (or a LF after CR)
			} else {
				// utils.SerialDbg("read irrelevant char", buf[i])
			}
		}

		if readFromCurrFrame+n < utils.SerialPacketSize {
			readFromCurrFrame += n
			driverutils.SerialPacket()
			print("rfcf ", n, " ", readFromCurrFrame)
			driverutils.SerialPacket()
		} else if readFromCurrFrame+n == utils.SerialPacketSize {
			driverutils.SerialAck()
			driverutils.SerialPacket()
			print("ack. rfcf=0")
			driverutils.SerialPacket()
			readFromCurrFrame = 0
		} else {
			driverutils.SerialPacket()
			print("tx error. in packet: ", readFromCurrFrame+n)
			driverutils.SerialPacket()
			return
		}
	}
}

func add(n uint32) {
	sum += n
	sumTerms++

	lcd.SetCursor(5, 0)
	utils.FormatUint32(sum, intBuf)
	lcd.Print(intBuf)
	lcd.SetCursor(5, 1)
	utils.FormatUint32(sumTerms, intBuf)
	lcd.Print(intBuf)
}

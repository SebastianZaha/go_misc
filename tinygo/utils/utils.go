package utils

import (
	"bufio"
	"io"
)

const (
	AsciiEOT = 4
	AsciiACK = 6
	AsciiLF  = 10
	AsciiCR  = 13
	AsciiUS  = 31 // Unit Separator, used for framing
	Ascii0   = 48
	Ascii9   = 57
)

// num = 1234, slice = [_,_,_,_]
// slice[0] = num / 10^3

var POW10 = []uint32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9}

// If the number does not fit, leftmost digit will be 'x'
func FormatUint32(num uint32, slice []byte) {
	var n = uint32(len(slice))
	var i uint32
	for i = 0; i < n; i++ {
		shift := POW10[n-i-1]
		//log.Println(num, shift, n, i)
		digit := num / shift
		if digit > 9 {
			digit = 'x' - 48
		}
		slice[i] = 48 + byte(digit)
		num = num % shift
	}
}

func Zero(slice []byte) {
	for i := 0; i < len(slice); i++ {
		slice[i] = 0
	}
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

const SerialPacketSize = 128

// alternate between waiting for ACK then writing a packet, until EOT
func InitSerialScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == AsciiACK {
				if i > 0 {
					return i, data[:i], nil
				} else {
					return i + 1, data[0:1], nil
				}
			} else if data[i] == AsciiUS {
				if i > 0 {
					return i + 1, data[:i], nil
				} else {
					// when first in buffer, simply skip
					return i + 1, nil, nil
				}
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		// There might be one final token to be delivered.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
		// but does not trigger an error to be returned from Scan itself.
		if len(data) == 0 {
			// illegal advance to stop the scan. no error and no token.
			// we do not want an empty last token when separator ends the buffer
			return 1, nil, nil
		} else {
			return 0, data, bufio.ErrFinalToken
		}
	}
	s.Split(split)
	return s
}

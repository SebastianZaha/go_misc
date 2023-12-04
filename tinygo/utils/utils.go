package utils

const (
	AsciiEOT = 4
	AsciiACK = 6
	AsciiLF  = 10
	AsciiSYN = 22
	AsciiCR  = 13
	AsciiUS  = 31 // Unit Separator
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

package main

import (
	"machine"
)

func main() {
	println("Reading until Ctrl-D")

	buf := make([]byte, 1024)
	sum := 0
	prevDigit := 0

	for {
		n, err := machine.Serial.Read(buf)
		if err != nil {
			println(err)
			return
		}
		if n > 0 {
			for i := 0; i < n; i++ {
				if buf[i] == 4 { // Ctrl-D, End of Input
					println("Sum: ", sum)
					return
				} else if buf[i] >= 48 && buf[i] < 58 {
					if prevDigit == 0 {
						prevDigit = int(buf[i])
						sum += 10 * (prevDigit - 48)
					} else {
						prevDigit = int(buf[i])
					}
					println("digit: ", buf[i], "; curr sum: ", sum)
				} else if buf[i] == 13 || buf[i] == 10 {
					sum += (prevDigit - 48)
					prevDigit = 0
					println("cr", sum)
				} else {
					// println("read irrelevant char", buf[i])
				}
			}
		}
	}
}

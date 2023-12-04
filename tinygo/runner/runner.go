package main

import (
	"bufio"
	"fmt"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"github.com/mattn/go-tty"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	err error

	serial *tty.TTY

	inputFileScanner *bufio.Scanner
)

func run() {
	serial, err = tty.OpenDevice("/dev/ttyACM0")
	if err != nil {
		fmt.Printf("Cannot open device: %+v", err)
		os.Exit(1)
	}

	scanner := initSerialScanner(serial.Input())
	for scanner.Scan() {
		token := scanner.Bytes()
		if len(token) == 0 {
			continue
		}
		if len(token) == 1 && token[len(token)-1] == utils.AsciiACK {
			writeInputPacket()
		} else {
			fmt.Printf("serial: %q\n", token)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanning serial input: %v", err)
		os.Exit(1)
	}
	fmt.Printf("finished scanning")
}

func writeInputPacket() {
	if inputFileScanner.Scan() {
		packet := inputFileScanner.Bytes()
		nwritten, err := serial.Output().Write(packet)
		utils.Must(err)
		if nwritten != len(packet) {
			fmt.Printf("Read %d from file, only wrote %d to serial", len(packet), nwritten)
			os.Exit(1)
		}
		fmt.Printf("Wrote from file to serial %d bytes\n%q\n\n", nwritten, packet)
		println()
		return
	}
	if err = inputFileScanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// if no token and no error we got to the end
	nwritten, err := serial.Output().Write([]byte{utils.AsciiEOT})
	utils.Must(err)
	if nwritten != 1 {
		fmt.Printf("Cannot write EOT byte, wrote %d bytes\n", nwritten)
		os.Exit(1)
	}
	fmt.Println("Wrote EOT to output")
}

// read the input file, split into utils.SerialPacketSize pieces
// CR and LF are squished into a single utils.AsciiUS because
// tty protocols mess up the CR/LF processing
func initFileScanner(f *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(f)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		token = make([]byte, utils.SerialPacketSize)
		dataIndex := 0
		tokenIndex := 0
		for ; tokenIndex < utils.SerialPacketSize && dataIndex < len(data); dataIndex++ {
			if data[dataIndex] == utils.AsciiCR || data[dataIndex] == utils.AsciiLF {
				if tokenIndex > 0 && token[tokenIndex-1] == utils.AsciiUS {
					// don't put a CR replacement if we already have one
				} else {
					token[tokenIndex] = utils.AsciiUS
					tokenIndex++
				}
			} else {
				token[tokenIndex] = data[dataIndex]
				tokenIndex++
			}
		}
		if tokenIndex == utils.SerialPacketSize {
			return dataIndex, token, nil
		} else if tokenIndex < utils.SerialPacketSize {
			if atEOF {
				// There is one final token to be delivered, which may be the empty string.
				// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
				// but does not trigger an error to be returned from Scan itself.
				return 0, token, bufio.ErrFinalToken
			} else {
				return 0, nil, nil
			}
		} else {
			log.Fatalf("file scanner logic error")
			return
		}
	}
	scanner.Split(split)
	return scanner
}

// alternate between waiting for ACK then writing a packet, until EOT
func initSerialScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == utils.AsciiACK {
				if i > 0 {
					return i, data[:i], nil
				} else {
					return i + 1, data[0:1], nil
				}
			} else if data[i] == utils.AsciiUS {
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

func main() {
	if len(os.Args) < 2 {
		fmt.Print(`Usage: 
	go run runner.go flash day1.go day1_1_ex.txt - will flash & run day1.go with provided test
	go run runner.go run day1.go day1_1_ex.txt - will feed day1_1.txt and monitor output 
`)
		return
	}
	goFile := os.Args[2]

	f, err := os.Open(os.Args[3])
	if err != nil {
		panic(err)
	}
	inputFileScanner = initFileScanner(f)

	if os.Args[1] == "flash" {
		args := []string{
			"flash",
			"-target=arduino",
			"-baudrate=9600",
			//"-gc=none",
			goFile,
		}
		log.Printf("Running:\n\ttinygo %s\n", args)
		cmd := exec.Command("tinygo", args...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	run()
}

// Testing using -monitor directly in the flash command, then redirecting stdio to it from our
// text test file. Does not seem to work because the tinygo programs interacts with the tty
// and ignores our stdio redirect somehow.
func testFlashAndMonitor() {
	cmd := exec.Command("tinygo", "flash", "-monitor", "-target=arduino", "-baudrate=9600", os.Args[1])

	outR, outW := io.Pipe()
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	cmd.Stdout = io.MultiWriter(os.Stdout, outW)
	cmd.Stderr = os.Stderr

	guard := "Reading until Ctrl-D"

	go func() {
		bufr := bufio.NewReader(outR)
		for {
			str, err := bufr.ReadString('\n')
			if err != nil {
				panic(err)
			}
			// stringCmpDebug(str, guard)
			str = strings.TrimSpace(str)
			if str == guard {
				f, err := os.Open(os.Args[2])
				if err != nil {
					panic(err)
				}
				n, err := io.Copy(inPipe, f)
				if err != nil {
					panic(err)
				}
				log.Printf("Copied %d bytes from %s.txt to program\n", n, os.Args[2])
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		panic(err)
	}
	if err = cmd.Wait(); err != nil {
		panic(err)
	}
}

func stringCmpDebug(str1, str2 string) {
	fmt.Print("Comparing strings: ")
	if len(str1) != len(str2) {
		fmt.Printf("lengths differ: %d vs %d\n", len(str1), len(str2))
		fmt.Printf("\t%v\n", []byte(str1))
		fmt.Printf("\t%v\n", []byte(str2))
		return
	}
	for i := 0; i < len(str1); i++ {
		if str1[i] != str2[i] {
			fmt.Printf("char %c at index %d differs from str2: %c\n", str1[i], i, str2[i])
			return
		}
	}
	fmt.Printf(" equal\n")
}

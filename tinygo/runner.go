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
	"time"
)

var (
	err error

	serial *tty.TTY

	inputFileScanner *bufio.Scanner
)

func run() {
	start := time.Now()
	serial, err = tty.OpenDevice("/dev/ttyACM0")
	if err != nil {
		fmt.Printf("Cannot open device: %+v", err)
		os.Exit(1)
	}

	scanner := utils.InitSerialScanner(serial.Input())
	for scanner.Scan() {
		token := scanner.Bytes()
		if len(token) == 0 {
			continue
		}
		if len(token) == 1 && token[len(token)-1] == utils.AsciiACK {
			writeInputPacket()
		} else if len(token) == 1 && token[len(token)-1] == utils.AsciiEOT {
			fmt.Printf("Finished in %.2fs\n", time.Since(start).Seconds())
			os.Exit(0)
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
		if len(packet) == 0 {
			fmt.Printf("Input file scanner should not return empty tokens")
			os.Exit(1)
		}
		n, err := serial.Output().Write([]byte{byte(len(packet))})
		if n != 1 || err != nil {
			fmt.Printf("Must be able to write 1 byte for the length of the packet. Wrote %d. Err: %v", n, err)
			os.Exit(1)
		}
		nwritten, err := serial.Output().Write(packet)
		utils.Must(err)
		if nwritten != len(packet) {
			fmt.Printf("Read %d from file, only wrote %d to serial", len(packet), nwritten)
			os.Exit(1)
		}
		// fmt.Printf("Wrote from file to serial %d bytes\n", nwritten)
		// fmt.Printf("\t%q\n\n", packet)
		return
	}
	if err = inputFileScanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// if no token and no error we got to the end
	nwritten, err := serial.Output().Write([]byte{1, utils.AsciiEOT})
	utils.Must(err)
	if nwritten != 2 {
		fmt.Printf("Cannot write EOT packet, wrote %d bytes\n", nwritten)
		os.Exit(1)
	}
	fmt.Println("Wrote EOT to output")
}

// Read the input file, split into utils.SerialPacketSize pieces
//
// CR and LF are squished and replaced by a single utils.AsciiUS because
// tty protocols mess with the CR/LF processing, seemingly adding or removing them
// in ways I do not understand.
func initFileScanner(f io.Reader) *bufio.Scanner {
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
				return 0, token[:tokenIndex], bufio.ErrFinalToken
			} else {
				// read more
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

func main() {
	if len(os.Args) < 2 {
		fmt.Print(`Usage: 
	go run runner.go flash day1_1/main.go day1_1/ex.txt - will flash & run with provided test
	go run runner.go run day1_1/ex.txt - will only feed text file and monitor output 
`)
		return
	}

	inFilePath := os.Args[2]
	if os.Args[1] == "flash" {
		inFilePath = os.Args[3]
		args := []string{
			"flash",
			"-target=arduino",
			"-baudrate=9600",
			// "-gc=leaking",
			// "-scheduler=tasks",
			os.Args[2],
		}
		log.Printf("Running:\n\ttinygo %s\n", args)
		cmd := exec.Command("tinygo", args...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	f, err := os.Open(inFilePath)
	if err != nil {
		panic(err)
	}
	inputFileScanner = initFileScanner(f)
	run()
}

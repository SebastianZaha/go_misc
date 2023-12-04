package main

import (
	"bufio"
	"fmt"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	f         *os.File
	err       error
	inFileBuf = make([]byte, utils.SerialPacketSize)

	ttyIn  *os.File
	ttyOut *os.File
)

func run(inputFile string) {
	ttyIn, err = os.Open("/dev/ttyACM0")
	if err != nil {
		fmt.Printf("Cannot open device: %+v", err)
		os.Exit(1)
	}
	ttyOut, err = os.OpenFile("/dev/ttyACM0", unix.O_WRONLY, 0)
	if err != nil {
		fmt.Printf("Cannot open device: %+v", err)
		os.Exit(1)
	}

	// alternate between waiting for ACK then writing a packet, until EOT
	scanner := bufio.NewScanner(ttyIn)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == utils.AsciiACK || data[i] == utils.AsciiUS {
				return i + 1, data[:i+1], nil
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
		// but does not trigger an error to be returned from Scan itself.
		return 0, data, bufio.ErrFinalToken
	}
	scanner.Split(split)
	for scanner.Scan() {
		token := scanner.Bytes()
		if len(token) == 1 {
			if token[0] == utils.AsciiACK {
				writeInputPacket()
			} else if token[0] != utils.AsciiUS {
				fmt.Printf("serial scanner got unexpected token separator %v", token[0])
				os.Exit(1)
			}
		} else if len(token) > 1 {
			fmt.Printf("serial: %q\n", token)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanning serial input: %v", err)
		os.Exit(1)
	}
}

func writeInputPacket() {
	n, err := f.Read(inFileBuf)
	if err != nil {
		fmt.Printf("cannot read from input file: %v", err)
		os.Exit(1)
	}
	if n > 0 {
		nwritten, err := ttyOut.Write(inFileBuf[:n])
		utils.Must(err)
		if nwritten != n {
			fmt.Printf("Read %d from file, only wrote %d to serial", n, nwritten)
			os.Exit(1)
		}
		fmt.Printf("Wrote from file to serial %d bytes\n%q\n\n", nwritten, inFileBuf[:n])
		println()
	} else {
		nwritten, err := ttyOut.Write([]byte{utils.AsciiEOT})
		utils.Must(err)
		if nwritten != 1 {
			fmt.Printf("Cannot write EOT byte, wrote %d bytes", nwritten)
			os.Exit(1)
		}
		fmt.Println("Wrote EOT to output")
	}
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
	inputFile := os.Args[3]
	f, err = os.Open(inputFile)
	if err != nil {
		panic(err)
	}

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

	run(inputFile)
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

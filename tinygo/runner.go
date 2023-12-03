package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-tty"
)

func run() {
	t, err := tty.OpenDevice("/dev/ttyACM0")
	if err != nil {
		log.Fatalf("Cannot open device: %+v", err)
	}
	guard := "Reading until Ctrl-D"

	bufr := bufio.NewReader(t.Input())
	for {
		str, err := bufr.ReadString('\n')
		if err != nil {
			log.Fatalf("Cannot read form tty output: %+v", err)
		}
		if strings.Contains(str, guard) {
			f, err := os.Open(os.Args[2] + ".txt")
			if err != nil {
				log.Fatal(err)
			}
			n, err := io.Copy(t.Output(), f)
			if err != nil {
				log.Fatalf("Error copying file to tty input: %v", err)
			}
			log.Printf("Copied %d bytes from %s.txt to program\n", n, os.Args[1])
		} else {
			print(str) // echo to tty
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(`Usage: 
	go run runner.go flash day1 - will flash & run day1.go
	go run runner.go run day1 - will feed day1_1.txt and monitor output 
`)
		return
	}

	if os.Args[1] == "flash" {
		cmd := exec.Command("tinygo", "flash", "-target=arduino", "-baudrate=9600", os.Args[2]+".go")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	run()
}

func testFlashAndMonitor() {
	cmd := exec.Command("tinygo", "flash", "-monitor", "-target=arduino", "-baudrate=9600", os.Args[2]+".go")

	outR, outW := io.Pipe()
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = io.MultiWriter(os.Stdout, outW)
	cmd.Stderr = os.Stderr

	guard := "Reading until Ctrl-D"

	go func() {
		bufr := bufio.NewReader(outR)
		for {
			str, err := bufr.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			// stringCmpDebug(str, guard)
			str = strings.TrimSpace(str)
			if str == guard {
				f, err := os.Open(os.Args[1] + ".txt")
				if err != nil {
					log.Fatal(err)
				}
				n, err := io.Copy(inPipe, f)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Copied %d bytes from %s.txt to program\n", n, os.Args[1])
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
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

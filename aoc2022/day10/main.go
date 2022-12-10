package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, _ := os.ReadFile("input.txt")
	var (
		reg    int = 1
		screen     = ""
	)

	lines := strings.Split(string(f), "\n")
	clk := 0
	lineIdx := 0
	incompleteInstruction := false
	screenLines := 0

	for {
		screenPos := clk % 40 // 0 to 39

		if screenPos == 0 {
			screen += "\n"
			screenLines++
			if screenLines == 7 {
				break
			}
		}

		if reg-1 == screenPos || reg == screenPos || reg+1 == screenPos {
			screen += "#"
		} else {
			screen += "."
		}

		// finish executing the "add" started on previous tick
		if incompleteInstruction {
			num, _ := strconv.ParseInt(lines[lineIdx][5:], 10, 64)
			reg += int(num)
			lineIdx++
			incompleteInstruction = false
		} else {
			// read new instruction
			if lines[lineIdx] == "noop" || len(lines[lineIdx]) == 0 {
				lineIdx++
			} else {
				incompleteInstruction = true
			}
		}

		clk += 1
	}
	fmt.Println(screen)
}

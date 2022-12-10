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
		cycles int64
		reg    int64 = 1
		screen       = ""
	)

	for _, l := range strings.Split(string(f), "\n") {
		if len(l) == 0 {
			continue
		}
		if l == "noop" {
			cycles += 1

			check := (cycles - 1) % 40
			if check == 0 {
				screen += "\n"
			}

			if reg-1 == check || reg == check || reg+1 == check {
				screen += "#"
			} else {
				screen += "."
			}
		} else {
			num, _ := strconv.ParseInt(l[5:], 10, 64)
			reg += num
			cycles += 2
			check := (cycles - 2) % 40
			check1 := (cycles - 1) % 40
			if check1 == 0 {
				screen += "\n"
			}
			if reg-num-1 == check || reg-num == check || reg-num+1 == check {
				screen += "#"
			} else {
				screen += "."
			}
			if check1 == 0 {
				screen += "\n"
			}
			if reg-num-1 == check1 || reg-num == check1 || reg-num+1 == check1 {
				screen += "#"
			} else {
				screen += "."
			}
		}
	}
	fmt.Println(screen)
}

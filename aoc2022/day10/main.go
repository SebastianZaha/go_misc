package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

func main() {
	t := time.Now()
	f, _ := os.ReadFile("input.txt")

	var (
		cycles      int64
		reg         int64 = 1
		sum         int64
		interesting = []int64{20, 60, 100, 140, 180, 220}
	)

	for _, l := range strings.Split(string(f), "\n") {
		if len(l) == 0 {
			continue
		}
		fmt.Printf("%d %s\n", cycles, l)

		if l == "noop" {
			cycles += 1

			if idx := slices.Index(interesting, cycles); idx > 0 {
				val := interesting[idx] * reg
				sum += val
				fmt.Printf("%d, %d, %d, %d, %d\n", interesting[idx], cycles, reg, val, sum)
			}
			continue
		}
		num, _ := strconv.ParseInt(l[5:], 10, 64)
		reg += num
		cycles += 2

		if idx := slices.Index(interesting, cycles-1); idx >= 0 {
			val := interesting[idx] * (reg - num)
			sum += val
			fmt.Printf("%d, %d, %d, %d, %d, %d\n", interesting[idx], cycles, reg, num, val, sum)
		} else if idx := slices.Index(interesting, cycles); idx >= 0 {
			val := interesting[idx] * (reg - num)
			sum += val
			fmt.Printf("%d, %d, %d, %d, %d\n", interesting[idx], cycles, reg, val, sum)
		}
	}

	fmt.Printf("Day10 1: %d. reg: %d, total cycles %d (in %v)\n", sum, reg, cycles, time.Since(t))
}

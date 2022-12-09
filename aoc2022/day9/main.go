package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var visited map[int]bool = map[int]bool{}

const INPUT = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

func moveSingleLink(direction byte, is []int, js []int, head int, tail int) {
	diffi := is[head] - is[tail]
	diffj := js[head] - js[tail]

	adiffi := diffi
	if adiffi < 0 {
		adiffi = 0 - adiffi
	}
	adiffj := diffj
	if adiffj < 0 {
		adiffj = 0 - adiffj
	}

	if adiffi == 2 {
		is[tail] += int(diffi / adiffi) // +-1
		if adiffj > 0 {
			js[tail] += int(diffj / adiffj) // +-1
		}
	} else if adiffj == 2 {
		js[tail] += int(diffj / adiffj)
		if adiffi > 0 {
			is[tail] += int(diffi / adiffi)
		}
	}

	//fmt.Printf("after move %c, is[head] %d, js[head] %d, is[tail] %d, js[tail] %d\n", direction, is[head], js[head], is[tail], js[tail])
}

func move(direction byte, is []int, js []int) {
	switch direction {
	case 'R':
		is[0] += 1
	case 'D':
		js[0] -= 1
	case 'L':
		is[0] -= 1
	case 'U':
		js[0] += 1
	}

	for i := 0; i < len(is)-1; i++ {
		moveSingleLink(direction, is, js, i, i+1)
	}
	fmt.Printf("%+v\n%+v\n\n", is, js)
	visit(is[len(is)-1], js[len(is)-1])
}

func visit(i, j int) {
	visited[i*100000+j] = true
}

func main() {
	// is := []int{1000, 1000}
	// js := []int{1000, 1000}
	is := []int{1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000}
	js := []int{1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000}

	f, _ := os.ReadFile("input.txt")
	for _, l := range strings.Split(string(f), "\n") {
		if len(l) == 0 {
			continue
		}
		num, _ := strconv.ParseInt(l[2:], 10, 64)
		for i := 0; i < int(num); i++ {
			move(l[0], is, js)
		}
	}
	fmt.Println(len(visited))

}

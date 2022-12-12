package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

var (
	h       int
	w       int
	si, sj  int
	ei, ej  int
	visited [][]int
	elev    [][]int
)

func printM(m [][]int) {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			fmt.Printf("%4d", m[i][j])
		}
		fmt.Println("")
	}
}

func tryGo(fromI, fromJ, toI, toJ, steps int) bool {
	// perimeter
	if toI < 0 || toJ < 0 || toI >= h || toJ >= w {
		return false
	}
	// already been here, with lower cost
	if visited[toI][toJ] > 0 && visited[toI][toJ] < visited[fromI][fromJ] {
		return false
	}
	// elev too high
	if elev[toI][toJ] > elev[fromI][fromJ]+1 {
		//fmt.Printf("elev too high %s %s\n", elev[toI][toJ], elev[fromI][fromJ])
		return false
	}

	visit(toI, toJ, steps)

	return true
}

var totalVisits = 0
var foundDestination = false

func visit(i int, j int, steps int) {
	totalVisits++
	if totalVisits > 10000000 {
		fmt.Println("total visits exceeded")
		return
	}
	if foundDestination {
		return
	}

	visited[i][j] = steps

	if i == ei && j == ej {
		foundDestination = true
		return
	}

	tryGo(i, j, i, j+1, steps+1) // right
	tryGo(i, j, i+1, j, steps+1) // down
	tryGo(i, j, i, j-1, steps+1) // up
	tryGo(i, j, i-1, j, steps+1) // left
}

func main() {

	/*	f := `Sabqponm
		abcryxxl
		accszExk
		acctuvwj
		abdefghi`*/
	f, _ := os.ReadFile("input.txt")
	f = bytes.TrimSpace(f)
	lines := strings.Split(string(f), "\n")

	h = len(lines)
	w = len(lines[0])

	visited = make([][]int, h)
	for i := range visited {
		visited[i] = make([]int, w)
	}
	elev = make([][]int, h)
	for i := range elev {
		elev[i] = make([]int, w)
	}

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if lines[i][j] == 'S' {
				si = i
				sj = j
				elev[i][j] = 0
			} else if lines[i][j] == 'E' {
				ei = i
				ej = j
				elev[i][j] = int('z' - 'a')
			} else {
				elev[i][j] = int(lines[i][j] - 'a')
			}
		}
	}

	fmt.Printf("Read matrix height %d, width %d, start at %d:%d end at %d:%d\n",
		h, w, si, sj, ei, ej)
	printM(elev)
	fmt.Println("")

	visit(si, sj, 0)

	printM(visited)
	fmt.Printf("Day12 lenght: %d\n", visited[ei][ej])
}

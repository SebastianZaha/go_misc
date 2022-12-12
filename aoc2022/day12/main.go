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
	ei, ej  int
	visited [][]int
	elev    [][]int
	stack   [][]int
)

func printM(m [][]int) {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			fmt.Printf("%4d", m[i][j])
		}
		fmt.Println("")
	}
}

func tryGo(fromI, fromJ, toI, toJ, steps int) {
	// perimeter
	if toI < 0 || toJ < 0 || toI >= h || toJ >= w {
		return
	}
	// already been here, with lower cost
	if visited[toI][toJ] > 0 && visited[toI][toJ] <= steps {
		return
	}
	// elev too high
	if elev[toI][toJ] > elev[fromI][fromJ]+1 {
		//fmt.Printf("elev too high %s %s\n", elev[toI][toJ], elev[fromI][fromJ])
		return
	}

	visited[toI][toJ] = steps
	stack = append(stack, []int{toI, toJ, steps})
}

func distFromStartingPoint(si, sj int) int {
	visited = make([][]int, h)
	for i := range visited {
		visited[i] = make([]int, w)
	}

	stack = [][]int{[]int{si, sj, 1}}
	visited[si][sj] = 1

	i := 0
	for {
		i++
		if i > 300000 {
			fmt.Printf("infinite loop: %d cycles\n", i)
			break
		}
		if len(stack) == 0 {
			break
		}
		pos := stack[0]
		stack = stack[1:]
		tryGo(pos[0], pos[1], pos[0], pos[1]+1, pos[2]+1) // right
		tryGo(pos[0], pos[1], pos[0]+1, pos[1], pos[2]+1) // down
		tryGo(pos[0], pos[1], pos[0], pos[1]-1, pos[2]+1) // up
		tryGo(pos[0], pos[1], pos[0]-1, pos[1], pos[2]+1) // left
	}

	// fmt.Printf("checked for %d cycles, found distance %d\n", i, visited[ei][ej])

	// -1 as we started with step 1 on the start position
	// (to distinguish from all the 0s that are default)
	return visited[ei][ej] - 1
}

func main() {

	/*f := `Sabqponm
	abcryxxl
	accszExk
	acctuvwj
	abdefghi`*/
	f, _ := os.ReadFile("input.txt")
	f = bytes.TrimSpace(f)
	lines := strings.Split(string(f), "\n")

	h = len(lines)
	w = len(lines[0])

	elev = make([][]int, h)
	for i := range elev {
		elev[i] = make([]int, w)
	}

	var possibleStarts [][]int = [][]int{}

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if lines[i][j] == 'S' || lines[i][j] == 'a' {
				elev[i][j] = 0
				possibleStarts = append(possibleStarts, []int{i, j})
			} else if lines[i][j] == 'E' {
				ei = i
				ej = j
				elev[i][j] = int('z' - 'a')
			} else {
				elev[i][j] = int(lines[i][j] - 'a')
			}
		}
	}

	fmt.Printf("Read matrix height %d, width %d, start at %d:%d end at %d:%d\n", h, w, ei, ej)
	printM(elev)
	fmt.Println("")

	min := 9999
	for _, pos := range possibleStarts {
		d := distFromStartingPoint(pos[0], pos[1])
		if d < min && d > 0 {
			min = d
		}
		// printM(visited)
	}

	fmt.Printf("Day12 min lenght: %d\n", min)
}

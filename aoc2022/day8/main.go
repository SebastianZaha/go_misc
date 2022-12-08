package main

import (
	"fmt"
	"os"
	"strings"
)

const EXAMPLE = `30373
25512
65332
33549
35390`

func main() {
	part1()
	part2()
}

func part2() {
	f, _ := os.ReadFile("input.txt")

	var (
		lines = strings.Split(string(f), "\n")
		rows  = len(lines)
		cols  = len(lines[0])
	)

	max := 0
	for i := 1; i < rows-1; i++ {
		for j := 1; j < cols-1; j++ {

			// left
			lscore := 0
			for k := j - 1; k >= 0; k-- {
				lscore++
				if lines[i][k] >= lines[i][j] {
					break
				}
			}
			// right
			rscore := 0
			for k := j + 1; k < cols; k++ {
				rscore++
				if lines[i][k] >= lines[i][j] {
					break
				}
			}
			// top
			tscore := 0
			for k := i - 1; k >= 0; k-- {
				tscore++
				if lines[k][j] >= lines[i][j] {
					break
				}
			}
			// bottom
			bscore := 0
			for k := i + 1; k < rows; k++ {
				bscore++
				if lines[k][j] >= lines[i][j] {
					break
				}
			}
			score := lscore * rscore * tscore * bscore
			if score > max {
				//fmt.Printf("%d,%d:%d,%d,%d,%d\n", i, j, lscore, rscore, tscore, bscore)
				max = score
			}
		}
	}
	fmt.Println(max)
}

// treat chars as ints. char "NINE" is ascii 57
const NINE = 57

func part1() {
	f, _ := os.ReadFile("input.txt")

	var (
		lines = strings.Split(string(f), "\n")
		rows  = len(lines)
		cols  = len(lines[0])
		// all edges are visible
		count = 4 * (len(lines) - 1)
		seen  = make([][]bool, rows)
	)

	for i := 0; i < rows; i++ {
		seen[i] = make([]bool, cols)
	}

	// for each of the 4 directions, do a single iteration (row OR col) iteration

	// left
	for i := 1; i < rows-1; i++ {
		max := lines[i][0]

		for j := 1; j < cols-1; j++ {
			if lines[i][j] > max {
				count += 1
				max = lines[i][j]
				seen[i][j] = true
				//fmt.Printf("%d,%d = %d\n", i, j, max)
				if max == NINE {
					break
				}
			}
		}
	}
	// right
	for i := 1; i < rows-1; i++ {
		max := lines[i][cols-1]

		for j := cols - 2; j > 0; j-- {
			// break early, we already did the rows in the opposite direction
			// also we don't need to check if seen when incrementing count lower
			if seen[i][j] {
				break
			}

			if lines[i][j] > max {
				count += 1
				max = lines[i][j]
				seen[i][j] = true
				//fmt.Printf("%d,%d = %d\n", i, j, max)
				if max == NINE {
					break
				}
			}
		}
	}
	// top
	for j := 1; j < cols-1; j++ {
		max := lines[0][j]

		for i := 1; i < rows-1; i++ {
			if lines[i][j] > max {
				if !seen[i][j] {
					count += 1
				}
				max = lines[i][j]
				seen[i][j] = true
				//fmt.Printf("%d,%d = %d\n", i, j, max)
				if max == NINE {
					break
				}
			}
		}
	}
	// bottom
	for j := 1; j < cols-1; j++ {
		max := lines[rows-1][j]

		for i := rows - 2; i > 0; i-- {

			if lines[i][j] > max {
				if !seen[i][j] {
					count += 1
				}
				max = lines[i][j]
				seen[i][j] = true
				//fmt.Printf("%d,%d = %d\n", i, j, max)
				if max == NINE {
					break
				}
			}
		}
	}

	fmt.Println(count)
}

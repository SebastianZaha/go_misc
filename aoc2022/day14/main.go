package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	cave       [][]int = make([][]int, 600)
	minY               = 999
	maxX, maxY         = 0, 0
)

func print() {
	for i := 0; i <= maxX; i++ {
		for j := minY; j <= maxY; j++ {
			if cave[i][j] == 0 {
				fmt.Printf(".")
			} else if cave[i][j] == 1 {
				fmt.Printf("#")
			} else {
				fmt.Printf("o")
			}
		}
		fmt.Println("")
	}
}

func main() {
	f, _ := os.ReadFile("input2.txt")
	f = bytes.TrimSpace(f)
	lines := strings.Split(string(f), "\n")

	for i := 0; i < len(cave); i++ {
		cave[i] = make([]int, 1000)
	}

	for i := 0; i < len(lines); i += 1 {
		prevX, prevY := 0, 0
		for j, coord := range strings.Split(lines[i], " -> ") {
			xy := strings.Split(coord, ",")
			x, err := strconv.Atoi(xy[1])
			if err != nil {
				fmt.Printf("panic parsing line %s, coord is not int %s\n", lines[i], coord)
				os.Exit(1)
			}
			y, err := strconv.Atoi(xy[0])
			if err != nil {
				fmt.Printf("panic parsing line %s, coord is not int %s\n", lines[i], coord)
				os.Exit(1)
			}

			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}

			if j != 0 {
				if prevX == x {
					if prevY > y {
						for k := y; k <= prevY; k++ {
							cave[x][k] = 1
						}
					} else {
						for k := prevY; k <= y; k++ {
							cave[x][k] = 1
						}
					}
				} else if prevX > x {
					for k := x; k <= prevX; k++ {
						cave[k][y] = 1
					}
				} else {
					for k := prevX; k <= x; k++ {
						cave[k][y] = 1
					}
				}
			}

			prevX = x
			prevY = y
		}
	}

	maxX += 2
	for i := 0; i < 1000; i++ {
		cave[maxX][i] = 1
	}

	sand := -1
	for {
		sand++
		if sand > 100000 {
			fmt.Println("algorithm error, loop too high")
			print()
			os.Exit(1)
		}

		x, y := 0, 500
		if cave[x][y] == 2 {
			fmt.Printf("filled up. stopping. (%d)\n", sand)
			return
		}
		for {
			if x == maxX {
				fmt.Printf("sand goes out the bottom, stopping. (%d)\n", sand)
				print()
				return
			} else if cave[x+1][y] == 0 {
				x++
			} else if y == 0 {
				fmt.Printf("sand goes out the left, stopping. (%d)\n", sand)
				print()
				return
			} else if cave[x+1][y-1] == 0 {
				x++
				y--
			} else if y == 999 {
				fmt.Printf("sand goes out the right, stopping. (%d)\n", sand)
				print()
				return
			} else if cave[x+1][y+1] == 0 {
				x++
				y++
			} else {
				cave[x][y] = 2
				break
			}
		}
	}
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
)

type node struct {
	isLeaf   bool
	val      byte
	parent   *node
	children []*node
}

func parse(line string) *node {
	root := &node{}
	current := root

	// first char is always the opening of the "root" node, so we skip it,
	// root is already initialized
	i := 0
	for i < len(line) {
		c := line[i]
		if c == '[' {
			newN := &node{parent: current}
			current.children = append(current.children, newN)
			current = newN
		} else if c == ']' {
			current = current.parent
		} else if c == ',' {
		} else if c == '1' {
			// the only 2 digit possible number is 10 apparently
			// easier to special case it here than to make a generic
			// variable-number-of-digits int parser
			if line[i+1] == '0' {
				current.children = append(current.children, &node{isLeaf: true, val: 10})
				i += 1
			} else {
				current.children = append(current.children, &node{isLeaf: true, val: 1})
			}
		} else if c >= '0' && c <= '9' {
			current.children = append(current.children, &node{isLeaf: true, val: c - '0'})
		} else {
			fmt.Printf("parser error: encountered unexpected char %c\n", c)
			os.Exit(1)
		}
		i += 1
	}

	return root
}

// -1, 0, 1
func cmp(n1, n2 *node) int8 {
	if n1.isLeaf {
		if n2.isLeaf {
			if n1.val < n2.val {
				return -1
			} else if n1.val == n2.val {
				return 0
			} else {
				return 1
			}
		} else {
			return cmp(&node{children: []*node{n1}}, n2)
		}
	} else {
		if n2.isLeaf {
			return cmp(n1, &node{children: []*node{n2}})
		} else {
			for i, c1 := range n1.children {
				// If the right list runs out of items first, the inputs are not in the right order.
				if i == len(n2.children) {
					return 1
				}
				childCmp := cmp(c1, n2.children[i])
				if childCmp != 0 {
					return childCmp
				}
			}
			if len(n1.children) == len(n2.children) {
				return 0
			} else {
				// If the left list runs out of items first, the inputs are in the right order.
				return -1
			}
		}
	}
}

func main() {
	f, _ := os.ReadFile("input2.txt")
	f = bytes.TrimSpace(f)
	lines := strings.Split(string(f), "\n")

	sum := 0
	for i := 0; i < len(lines); i += 3 {
		v1 := parse(lines[i])
		v2 := parse(lines[i+1])
		cmpVal := cmp(v1, v2)
		if cmpVal == -1 {
			sum += i/3 + 1
		}
	}

	fmt.Printf("Day13.1 sum %d\n", sum)

	divider1 := parse("[[2]]")
	divider2 := parse("[[6]]")
	parsedLines := []*node{divider1, divider2}

	for _, l := range lines {
		if len(l) != 0 {
			parsedLines = append(parsedLines, parse(l))
		}
	}

	sort.Slice(parsedLines, func(i, j int) bool {
		return cmp(parsedLines[i], parsedLines[j]) == -1
	})

	prod := 1
	for i, l := range parsedLines {
		if l == divider1 || l == divider2 {
			prod *= i + 1
		}
	}

	fmt.Printf("Day13.2 prod %d\n", prod)
}

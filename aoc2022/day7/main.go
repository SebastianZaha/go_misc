package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type entry struct {
	name     string
	size     int64
	parent   *entry
	children map[string]*entry
}

func newEntry(name string, size int64, parent *entry) *entry {
	c := &entry{name: name, size: size, parent: parent, children: map[string]*entry{}}
	if parent != nil {
		parent.children[name] = c
	}
	return c
}

func walkDepthFirst(e *entry, cb func(*entry)) {
	for _, v := range e.children {
		walkDepthFirst(v, cb)
	}
	cb(e)
}

func main() {
	f, _ := ioutil.ReadFile("input.txt")

	root := newEntry("/", 0, nil)
	curr := root

	for _, line := range strings.Split(string(f), "\n") {
		if line == "" || line == "$ cd /" || line == "$ ls" {
			continue
		} else if line[0:5] == "$ cd " {
			name := line[5:]
			if name == ".." {
				curr = curr.parent
			} else if c, ok := curr.children[name]; ok {
				// fmt.Printf("current dir is now %s\n", c.name)
				curr = c
			} else {
				fmt.Printf("Cannot find dir to cd into: %s\n", name)
				os.Exit(1)
			}
		} else if line[0:4] == "dir " {
			// fmt.Printf("adding dir %s to %s\n", line[4:], curr.name)
			newEntry(line[4:], 0, curr)
		} else {
			parts := strings.Split(line, " ")
			size, _ := strconv.ParseInt(parts[0], 10, 64)
			// fmt.Printf("adding file %s, size %d to %s\n", parts[1], size, curr.name)
			newEntry(parts[1], size, curr)
		}
	}

	var totalPart1 int64
	walkDepthFirst(root, func(e *entry) {
		if e.parent != nil {
			e.parent.size += e.size
		}
		if len(e.children) > 0 && e.size < 100_000 {
			totalPart1 += e.size
		}
		// fmt.Printf("%s %d\n", e.name, e.size)
	})
	fmt.Printf("part1 total: %d\n", totalPart1)

	// fs has 70_000_000. 30_000_000 needs to be empty after deletion
	// presumably now 30_000_000 + root.size > 70_000_000
	needToFree := root.size + 30_000_000 - 70_000_000
	minToDel := root.size

	walkDepthFirst(root, func(e *entry) {
		if len(e.children) != 0 && e.size > needToFree && e.size < minToDel {
			minToDel = e.size
		}
	})

	fmt.Printf("needed: %d, minToDel: %d\n", needToFree, minToDel)
}

package main

import (
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"os"
	"slices"
	"testing"
)

func TestInputFileScanner(t *testing.T) {
	bs, err := os.ReadFile("day3_1/ex.txt")
	// change the separators manually as would be done by the scanner
	for i := 0; i < len(bs); i++ {
		if bs[i] == 10 {
			bs[i] = utils.AsciiUS
		}
	}

	if err != nil {
		t.Fatalf("cannot read test file")
	}
	f, err := os.Open("day3_1/ex.txt")
	if err != nil {
		t.Fatalf("cannot open test file")
	}
	s := initFileScanner(f)

	if !s.Scan() {
		t.Fatalf("should return at least a token")
	}
	if 0 != slices.Compare(s.Bytes(), bs) {
		t.Fatalf("The file contents should be entirely returned as a token. \nGot: %+v\nExp: %+v", s.Bytes(), bs)
	}

	if s.Scan() {
		t.Fatalf("should not be able to scan, entire file fits in a token. Got: bs: %+v \n err: %v", s.Bytes(), s.Err())
	}
}

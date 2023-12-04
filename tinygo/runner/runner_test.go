package main

import (
	"bytes"
	"fmt"
	"github.com/SebastianZaha/go_misc/tinygo/utils"
	"slices"
	"testing"
)

func TestSerialScanner(t *testing.T) {
	clean := []byte("no separators")

	tests := []struct {
		name     string
		in       []byte
		expected [][]byte
	}{
		{
			name:     "single",
			in:       clean,
			expected: [][]byte{clean},
		}, {
			name:     "ending with separator",
			in:       append(clean, utils.AsciiUS),
			expected: [][]byte{clean},
		}, {
			name:     "two with separator",
			in:       append(append(clean, utils.AsciiUS), clean...),
			expected: [][]byte{clean, clean},
		}, {
			name:     "ACK by itself",
			in:       []byte{utils.AsciiACK},
			expected: [][]byte{{utils.AsciiACK}},
		}, {
			name:     "clean ACK clean",
			in:       append(append(clean, utils.AsciiACK), clean...),
			expected: [][]byte{clean, {utils.AsciiACK}, clean},
		}, {
			name:     "clean ACK US clean",
			in:       append(append(clean, utils.AsciiACK, utils.AsciiUS), clean...),
			expected: [][]byte{clean, {utils.AsciiACK}, clean},
		},
	}

	for _, tt := range tests {
		fmt.Printf("starting test %s\n", tt.name)

		r := bytes.NewReader(tt.in)
		s := initSerialScanner(r)
		for i, expected := range tt.expected {
			if !s.Scan() {
				t.Errorf("%s: Expected to scan token: %v", tt.name, s.Err())
			}
			if c := slices.Compare(s.Bytes(), expected); c != 0 {
				t.Errorf("%s: Expected token %d to be %v, was %v", tt.name, i, expected, s.Bytes())
			}
		}
		if s.Scan() {
			t.Errorf("%s: Should not have any more tokens. Got %v", tt.name, s.Bytes())
		}
	}
}

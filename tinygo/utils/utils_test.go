package utils

import (
	"bytes"
	"fmt"
	"slices"
	"testing"
)

func TestFormatUint32(t *testing.T) {
	tests := []struct {
		num      uint32
		expected string
	}{
		{3, "000003"},
		{22, "000022"},
		{123456, "123456"},
		{1234567, "x34567"},
	}

	intSlice := make([]byte, 6)
	for _, tt := range tests {
		FormatUint32(tt.num, intSlice)
		if string(intSlice) != tt.expected {
			t.Errorf("Expected %s, got %+v (%s)", tt.expected, intSlice, string(intSlice))
		}
	}
}

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
			in:       append(clean, AsciiUS),
			expected: [][]byte{clean},
		}, {
			name:     "two with separator",
			in:       append(append(clean, AsciiUS), clean...),
			expected: [][]byte{clean, clean},
		}, {
			name:     "ACK by itself",
			in:       []byte{AsciiACK},
			expected: [][]byte{{AsciiACK}},
		}, {
			name:     "clean ACK clean",
			in:       append(append(clean, AsciiACK), clean...),
			expected: [][]byte{clean, {AsciiACK}, clean},
		}, {
			name:     "clean ACK US clean",
			in:       append(append(clean, AsciiACK, AsciiUS), clean...),
			expected: [][]byte{clean, {AsciiACK}, clean},
		},
	}

	for _, tt := range tests {
		fmt.Printf("starting test %s\n", tt.name)

		r := bytes.NewReader(tt.in)
		s := InitSerialScanner(r)
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

package utils

import (
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

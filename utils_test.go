package main

import (
	"testing"
)

func TestToUint8(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected uint8
	}{
		{int32(60), 60},
		{int64(127), 127},
		{float32(100.5), 100},
		{float64(64.7), 64},
		{int(42), 42},
		{uint(88), 88},
		{256, 0},       // Should handle overflow
		{-1, 0},        // Should handle negative
		{"invalid", 0}, // Should handle non-numeric
		{nil, 0},       // Should handle nil
	}

	for _, tt := range tests {
		result := toUint8(tt.input)
		if result != tt.expected {
			t.Errorf("toUint8(%v) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

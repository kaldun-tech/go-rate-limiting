package practice

import (
	"testing"
)

func TestClosestToZero_Empty(t *testing.T) {
	got := ClosestToZero([]float64{})
	if got != 0 {
		t.Errorf("ClosestToZero([]) = %v, want 0", got)
	}
}

func TestClosestToZero_SingleValue(t *testing.T) {
	tests := []struct {
		input    []float64
		expected float64
	}{
		{[]float64{5.0}, 5.0},
		{[]float64{-3.0}, -3.0},
		{[]float64{0.0}, 0.0},
	}

	for _, tc := range tests {
		got := ClosestToZero(tc.input)
		if got != tc.expected {
			t.Errorf("ClosestToZero(%v) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestClosestToZero_MultipleValues(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"positive values", []float64{1.0, 2.0, 3.0}, 1.0},
		{"negative values", []float64{-1.0, -2.0, -3.0}, -1.0},
		{"mixed values", []float64{-5.0, 3.0, -1.0, 2.0}, -1.0},
		{"sample from problem", []float64{7.0, -10.0, 13.0, 8.0, 4.0, -7.2, -1.7, 6.0, -4.0}, -1.7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ClosestToZero(tc.input)
			if got != tc.expected {
				t.Errorf("ClosestToZero(%v) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestClosestToZero_TieBreaker(t *testing.T) {
	// When equally close, prefer positive
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"5 and -5", []float64{-5.0, 5.0}, 5.0},
		{"reversed order", []float64{5.0, -5.0}, 5.0},
		{"1 and -1 with others", []float64{-10.0, -1.0, 1.0, 10.0}, 1.0},
		{"small decimals", []float64{-0.5, 0.5}, 0.5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ClosestToZero(tc.input)
			if got != tc.expected {
				t.Errorf("ClosestToZero(%v) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestClosestToZero_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"contains zero", []float64{-5.0, 0.0, 5.0}, 0.0},
		{"extreme cold", []float64{-273.0, -100.0}, -100.0},
		{"extreme hot", []float64{5526.0, 1000.0}, 1000.0},
		{"very close decimals", []float64{0.1, -0.2, 0.15}, 0.1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ClosestToZero(tc.input)
			if got != tc.expected {
				t.Errorf("ClosestToZero(%v) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

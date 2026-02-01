package practice

import (
	"testing"
)

func TestComputeTargetHeights(t *testing.T) {
	tests := []struct {
		name     string
		boxes    []int
		expected []int
	}{
		{"even distribution", []int{2, 2, 2}, []int{2, 2, 2}},
		{"needs leveling", []int{4, 1, 1}, []int{2, 2, 2}},
		{"with remainder", []int{3, 2, 2}, []int{3, 2, 2}}, // 7 boxes, 3 stacks -> [3,2,2]
		{"two stacks even", []int{3, 3}, []int{3, 3}},
		{"two stacks odd", []int{4, 1}, []int{3, 2}}, // 5 boxes, 2 stacks -> [3,2]
		{"all on one", []int{6, 0, 0}, []int{2, 2, 2}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := computeTargetHeights(tc.boxes)
			if len(got) != len(tc.expected) {
				t.Errorf("computeTargetHeights(%v) length = %d, want %d", tc.boxes, len(got), len(tc.expected))
				return
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("computeTargetHeights(%v) = %v, want %v", tc.boxes, got, tc.expected)
					return
				}
			}
		})
	}
}

func TestIsLeveled(t *testing.T) {
	tests := []struct {
		name     string
		boxes    []int
		targets  []int
		expected bool
	}{
		{"already leveled", []int{2, 2, 2}, []int{2, 2, 2}, true},
		{"not leveled", []int{4, 1, 1}, []int{2, 2, 2}, false},
		{"leveled with remainder", []int{3, 2, 2}, []int{3, 2, 2}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isLeveled(tc.boxes, tc.targets)
			if got != tc.expected {
				t.Errorf("isLeveled(%v, %v) = %v, want %v", tc.boxes, tc.targets, got, tc.expected)
			}
		})
	}
}

func TestFindSource(t *testing.T) {
	tests := []struct {
		name     string
		boxes    []int
		targets  []int
		expected int
	}{
		{"first stack is source", []int{4, 1, 1}, []int{2, 2, 2}, 0},
		{"middle stack is source", []int{2, 4, 0}, []int{2, 2, 2}, 1},
		{"no source needed", []int{2, 2, 2}, []int{2, 2, 2}, -1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := findSource(tc.boxes, tc.targets)
			if got != tc.expected {
				t.Errorf("findSource(%v, %v) = %v, want %v", tc.boxes, tc.targets, got, tc.expected)
			}
		})
	}
}

func TestFindTarget(t *testing.T) {
	tests := []struct {
		name     string
		boxes    []int
		targets  []int
		expected int
	}{
		{"second stack needs boxes", []int{4, 1, 1}, []int{2, 2, 2}, 1},
		{"first stack needs boxes", []int{0, 3, 3}, []int{2, 2, 2}, 0},
		{"no target needed", []int{2, 2, 2}, []int{2, 2, 2}, -1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := findTarget(tc.boxes, tc.targets)
			if got != tc.expected {
				t.Errorf("findTarget(%v, %v) = %v, want %v", tc.boxes, tc.targets, got, tc.expected)
			}
		})
	}
}

// Simulation test - runs Solve repeatedly to verify it reaches the goal
func TestSolve_Simulation(t *testing.T) {
	tests := []struct {
		name         string
		initialBoxes []int
		initialClaw  int
	}{
		{"simple case", []int{4, 1, 1}, 0},
		{"already leveled", []int{2, 2, 2}, 0},
		{"all on left", []int{6, 0, 0}, 0},
		{"all on right", []int{0, 0, 6}, 0},
		{"with remainder", []int{5, 1, 1}, 0},
		{"two stacks", []int{4, 0}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			boxes := make([]int, len(tc.initialBoxes))
			copy(boxes, tc.initialBoxes)
			clawPos := tc.initialClaw
			boxInClaw := false

			targets := computeTargetHeights(tc.initialBoxes)
			if targets == nil {
				t.Skip("computeTargetHeights not implemented")
			}

			for turn := 0; turn < 200; turn++ {
				if isLeveled(boxes, targets) {
					return // Success
				}

				cmd := Solve(clawPos, boxes, boxInClaw)

				// Execute command
				switch cmd {
				case Right:
					if clawPos < len(boxes)-1 {
						clawPos++
					} else {
						t.Fatalf("Invalid RIGHT at position %d", clawPos)
					}
				case Left:
					if clawPos > 0 {
						clawPos--
					} else {
						t.Fatalf("Invalid LEFT at position %d", clawPos)
					}
				case Pick:
					if boxInClaw {
						t.Fatal("Cannot PICK while holding a box")
					}
					if boxes[clawPos] == 0 {
						t.Fatalf("Cannot PICK from empty stack %d", clawPos)
					}
					boxes[clawPos]--
					boxInClaw = true
				case Place:
					if !boxInClaw {
						t.Fatal("Cannot PLACE without holding a box")
					}
					boxes[clawPos]++
					boxInClaw = false
				default:
					t.Fatalf("Unknown command: %q", cmd)
				}
			}

			t.Errorf("Failed to level boxes within 200 turns. Final state: %v", boxes)
		})
	}
}

package practice

// Robotic Leveling
//
// Goal: Rearrange boxes in the factory to form stacks of equal height.
//
// Rules:
// You work in an automated factory with a simple robotic arm that can move
// boxes around. The arm can pick a box from a stack and place it on another.
//
// Your objective is to rearrange the stacks to have an equal number of boxes
// on each stack. If not possible, excess boxes must be stacked from left to right.
//
// Example: 7 boxes across 3 stacks -> [3, 2, 2] (leftmost gets the extra)
//
// Available Commands:
// - "RIGHT": arm moves one stack to the right
// - "LEFT": arm moves one stack to the left
// - "PICK": arm grabs a box from the stack below
// - "PLACE": arm places a box onto the stack below
//
// Victory: All stacks are leveled (equal height, extras distributed left to right)
// Lose: Invalid command or not solved within 200 turns
//
// Constraints:
// - 2 <= number of stacks <= 8
// - 1 <= number of boxes <= 16

// Command constants
const (
	Right = "RIGHT"
	Left  = "LEFT"
	Pick  = "PICK"
	Place = "PLACE"
)

// Solve returns the next command to execute given the current state.
// Caller is responsible for checking whether the puzzle is already solved and in a valid state
// Parameters:
//   - clawPos: index of the stack the arm is currently above (0-indexed)
//   - boxes: slice of integers representing the height of each stack
//   - boxInClaw: true if the arm is currently holding a box
//
// Returns: one of "RIGHT", "LEFT", "PICK", or "PLACE", ""
// The empty string indicates caller error due to either a solved puzzle or inconsistent state
// O(n) for n length of boxes
func Solve(clawPos int, boxes []int, boxInClaw bool) string {
	// Hints:
	// - First, calculate the target height for each stack
	//   totalBoxes / numStacks = base height
	//   totalBoxes % numStacks = number of stacks that get +1 (from left)
	targets := computeTargetHeights(boxes, boxInClaw)
	// Identify which stack has too many boxes - source
	source := findSource(boxes, targets)
	// Identify which stack needs more boxes - destination
	dest := findTarget(boxes, targets)

	// Strategy: move to a source, pick, move to a target, place
	if boxInClaw && 0 <= dest {
		// Move to the destination and place
		if clawPos == dest {
			return Place
		} else if clawPos < dest {
			return Right
		} else {
			return Left
		}
	} else if !boxInClaw && 0 <= source {
		// Move to the source and pick
		if clawPos < source {
			return Right
		} else if clawPos == source {
			return Pick
		} else {
			return Left
		}
	}
	// Error: already solved or invalid state
	return ""
}

// Helper: compute target heights for a leveled configuration and claw state
// Returns a slice where targetHeights[i] is the desired height for stack i
// O(n) for n length of boxes
func computeTargetHeights(boxes []int, boxInClaw bool) []int {
	// Total boxes divided among stacks, extras go left to right
	sum := 0
	n := len(boxes)
	for _, b := range boxes {
		sum += b
	}
	// Add one if there is a box in the claw
	if boxInClaw {
		sum++
	}

	// Compute integer mean as lower bound and remainder
	avg := sum / n
	rem := sum % n

	// Each stack gets the lower bound of boxes, plus one while there is remainder
	targetHeights := make([]int, n)
	for i := range n {
		targetHeights[i] = avg
		if 0 < rem {
			targetHeights[i]++
			rem--
		}
	}
	return targetHeights
}

// Helper: check if current configuration matches target
// Assumes the inputs are of matched length
// O(n) for n length of boxes
func isLeveled(boxes []int, targets []int) bool {
	for i, b := range boxes {
		if b != targets[i] {
			// Mismatch
			return false
		}
	}
	// All match
	return true
}

// Helper: find the leftmost stack that has more boxes than target
// Returns index of a stack with boxes[i] > targets[i], or -1 if none
// O(n) for n length of boxes
func findSource(boxes []int, targets []int) int {
	for i, b := range boxes {
		t := targets[i]
		if t < b {
			return i
		}
	}
	// None found
	return -1
}

// Helper: find the leftmost stack that needs more boxes
// Returns index of a stack with boxes[i] < targets[i], or -1 if none
// O(n) for n length of boxes
func findTarget(boxes []int, targets []int) int {
	for i, b := range boxes {
		t := targets[i]
		if b < t {
			return i
		}
	}
	// None found
	return -1
}

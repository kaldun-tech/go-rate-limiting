package practice

import "math"

// Find the Temperature
//
// Your company builds temperature captors for freezers. These captors record
// temperature periodically and put the last values in an array. You have to
// develop the algorithm displaying the unique temperature that is supposed to
// sum up these values.
//
// You know the captors are not reliable at all, so you decide to display the
// most expected temperature among the ones in the array, which is the one
// closest to zero.
//
// Rules:
// - The array ts is always defined (never nil)
// - ts can be empty, in that case, return 0
// - If two numbers are equally close to zero, return the positive one
//   (e.g., if ts contains -5 and 5, return 5)
// - Temperatures range from -273.0 to 5526.0

// ClosestToZero returns the temperature closest to zero from the given slice.
// If the slice is empty, returns 0.
// If two temperatures are equally close, returns the positive one.
// O(n) for n length of input ts
func ClosestToZero(ts []float64) float64 {
	// Hints:
	// - Handle empty slice case first
	// - Track the closest temperature seen so far
	// - Use math.Abs() for comparing distances to zero
	// - Remember the tie-breaker rule: prefer positive values

	if len(ts) == 0 {
		// Empty -> return 0
		return 0
	}

	var closestToZero float64 = ts[0]
	for _, next := range ts {
		absClosest := math.Abs(closestToZero)
		absNext := math.Abs(next)
		if absNext == 0 {
			return 0
		}
		if absNext < absClosest || (absNext == absClosest && 0 < next && closestToZero < 0) {
			// Update closest for closer absolute or equal positive next value
			closestToZero = next
		}
	}

	return closestToZero
}

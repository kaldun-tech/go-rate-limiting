package datastructures

// MinHeap implements a binary min heap
// https://en.wikipedia.org/wiki/Binary_heap
// Common interview uses:
// - Find K largest/smallest elements
// - Merge K sorted lists
// - Priority queue implementation
// - Median from data stream
type MinHeap struct {
	items []int
}

// NewMinHeap creates an empty min heap
// Min heap: parent value must be less than child
func NewMinHeap() *MinHeap {
	return &MinHeap{
		items: make([]int, 0),
	}
}

// Push adds an element to the heap
// Time: O(log n)
func (h *MinHeap) Push(val int) {
	// 1. Add to end of array
	i := len(h.items)
	h.items = append(h.items, val)
	// 2. Bubble up to maintain heap property:
	for 0 < i && val < h.items[h.parent(i)] {
		p := h.parent(i)
		h.swap(i, p)
		i = p
	}
}

// Pop removes and returns the minimum element
// Returns: value, isValid (false for empty heap)
// Time: O(log n)
func (h *MinHeap) Pop() (int, bool) {
	// Handle empty case
	if h.IsEmpty() {
		return 0, false
	}

	// 1. Save root value and final value
	rootVal := h.items[0]
	fVal := h.items[len(h.items)-1]

	// 2. Move last element to root and shrink the slice
	h.items[0] = fVal
	h.items = h.items[:len(h.items)-1]

	// 3. Bubble down to maintain heap property
	for i := 0; ; {
		smallest := i
		l, r := h.leftChild(i), h.rightChild(i)
		n := len(h.items)

		// Check if children exist before accessing them (will panic on out-of-bounds)
		// Swap with the smaller child to guarantee the new parent is ≤ both children.
		// When you promote a child to parent, it becomes the parent of its sibling too.
		// If you promote the larger child, it's now parent of a smaller sibling → violation
		if l < n && h.items[l] < h.items[smallest] {
			smallest = l
		}
		if r < n && h.items[r] < h.items[smallest] {
			smallest = r
		}
		if smallest == i {
			break
		}
		h.swap(i, smallest)
		i = smallest
	}

	return rootVal, true
}

// Peek returns the minimum element without removing it
// Returns: value, isValid (false for empty heap)
// Time: O(1)
func (h *MinHeap) Peek() (int, bool) {
	// Handle empty case
	if h.IsEmpty() {
		return 0, false
	}
	return h.items[0], true
}

// Size returns the number of elements in the heap
func (h *MinHeap) Size() int {
	return len(h.items)
}

// Returns whether heap is empty
func (h *MinHeap) IsEmpty() bool {
	return h.Size() == 0
}

// Helper functions for heap operations

func (h *MinHeap) parent(i int) int {
	return (i - 1) / 2
}

func (h *MinHeap) leftChild(i int) int {
	return 2*i + 1
}

func (h *MinHeap) rightChild(i int) int {
	return 2*i + 2
}

func (h *MinHeap) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

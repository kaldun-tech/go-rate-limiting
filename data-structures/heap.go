package datastructures

// MinHeap implements a binary min heap
// Common interview uses:
// - Find K largest/smallest elements
// - Merge K sorted lists
// - Priority queue implementation
// - Median from data stream
type MinHeap struct {
	items []int
}

// NewMinHeap creates an empty min heap
func NewMinHeap() *MinHeap {
	return &MinHeap{
		items: make([]int, 0),
	}
}

// Push adds an element to the heap
// Time: O(log n)
func (h *MinHeap) Push(val int) {
	// TODO: Implement
	// 1. Add to end of array
	// 2. Bubble up to maintain heap property
}

// Pop removes and returns the minimum element
// Time: O(log n)
func (h *MinHeap) Pop() (int, bool) {
	// TODO: Implement
	// 1. Save root value
	// 2. Move last element to root
	// 3. Bubble down to maintain heap property
	return 0, false
}

// Peek returns the minimum element without removing it
// Time: O(1)
func (h *MinHeap) Peek() (int, bool) {
	// TODO: Implement
	return 0, false
}

// Size returns the number of elements in the heap
func (h *MinHeap) Size() int {
	return len(h.items)
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

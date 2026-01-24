package datastructures

// baseHeap contains shared fields and helper methods for heap implementations
type baseHeap struct {
	items []int
}

func (h *baseHeap) parent(i int) int {
	return (i - 1) / 2
}

func (h *baseHeap) leftChild(i int) int {
	return 2*i + 1
}

func (h *baseHeap) rightChild(i int) int {
	return 2*i + 2
}

func (h *baseHeap) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

// Size returns the number of elements in the heap
func (h *baseHeap) Size() int {
	return len(h.items)
}

// IsEmpty returns whether heap is empty
func (h *baseHeap) IsEmpty() bool {
	return len(h.items) == 0
}

// baseHeap has a generic siftDown that relies on a comparator
func (h *baseHeap) siftDown(i int, less func(i, j int) bool) {
	for {
		target := i
		l, r := h.leftChild(i), h.rightChild(i)
		n := len(h.items)

		// less returns true if items[i] should be higher in the heap than items[j]
		if l < n && less(l, target) {
			target = l
		}
		if r < n && less(r, target) {
			target = r
		}
		if target == i {
			break
		}
		h.swap(i, target)
		i = target
	}
}

// Peek returns the minimum element without removing it
// Returns: value, isValid (false for empty heap)
// Time: O(1)
func (h *baseHeap) Peek() (int, bool) {
	// Handle empty case
	if h.IsEmpty() {
		return 0, false
	}
	return h.items[0], true
}

// MinHeap implements a binary min heap
// https://en.wikipedia.org/wiki/Binary_heap
// Common interview uses:
// - Find K largest/smallest elements
// - Merge K sorted lists
// - Priority queue implementation
// - Median from data stream
type MinHeap struct {
	baseHeap
}

// NewMinHeap creates an empty min heap
// Min heap: parent value must be less than child
func NewMinHeap() *MinHeap {
	return &MinHeap{
		baseHeap: baseHeap{items: make([]int, 0)},
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
	h.siftDown(0)

	return rootVal, true
}

// Peek returns the minimum element without removing it
// Returns: value, isValid (false for empty heap)
// Time: O(1)
func (h *MinHeap) Peek() (int, bool) {
	return h.baseHeap.Peek()
}

// siftDown pushes the node at index i down the heap until heap property is restored
// Time: O(log n)
func (h *MinHeap) siftDown(i int) {
	// Alternative: rely on baseheap
	// ltcmp := func(i, j int) bool { return h.items[i] < h.items[j] }
	// h.baseheap.siftDdown(i, ltcmp)
	for {
		smallest := i
		l, r := h.leftChild(i), h.rightChild(i)
		n := len(h.items)

		// Check if children exist before accessing them
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
}

// Heapify converts an unordered slice into a valid min heap in-place
// Input slice is copied into this heap's items
// This is more efficient than inserting elements one by one
// Time: O(n) - not O(n log n)!
func (h *MinHeap) Heapify(items []int) {
	// I choose to copy data into this heap's items. ALternative is to take ownership of the array directly.
	h.items = make([]int, len(items))
	copy(h.items, items)

	// Hint: Start from last non-leaf node and sift down each node
	// Meaning pushing a node down the min heap until its smaller than both children
	// Last non-leaf is at index (n/2 - 1)
	for i := len(items)/2 - 1; 0 <= i; i-- {
		h.siftDown(i)
	}
}

// MaxHeap implements a binary max heap
type MaxHeap struct {
	baseHeap
}

// NewMaxHeap creates an empty max heap
func NewMaxHeap() *MaxHeap {
	return &MaxHeap{
		baseHeap: baseHeap{items: make([]int, 0)},
	}
}

// Push adds an element to the heap
// Time: O(log n)
func (h *MaxHeap) Push(val int) {
	// 1. Add to end of array
	i := len(h.items)
	h.items = append(h.items, val)
	// 2. Bubble up to maintain heap property:
	for 0 < i && h.items[h.parent(i)] < val {
		p := h.parent(i)
		h.swap(i, p)
		i = p
	}
}

// Pop removes and returns the maximum element
// Time: O(log n)
func (h *MaxHeap) Pop() (int, bool) {
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
	h.siftDown(0)

	return rootVal, true
}

// Peek returns the maximum element without removing it
// Time: O(1)
func (h *MaxHeap) Peek() (int, bool) {
	return h.baseHeap.Peek()
}

// siftDown pushes the node at index i down the heap until heap property is restored
// Time: O(log n)
func (h *MaxHeap) siftDown(i int) {
	// Alternative: rely on baseheap
	// gtcmp := func(i, j int) bool { return h.items[i] > h.items[j] }
	// h.baseheap.siftDdown(i, gtcmp)
	for {
		greatest := i
		l, r := h.leftChild(i), h.rightChild(i)
		n := len(h.items)

		// Check if children exist before accessing them
		// Swap with the larger child in this case to guarantee the new parent is >= both children.
		// When you promote a child to parent, it becomes the parent of its sibling too.
		// If you promote the smaller child, it's now parent of a smaller sibling → violation
		if l < n && h.items[greatest] < h.items[l] {
			greatest = l
		}
		if r < n && h.items[greatest] < h.items[r] {
			greatest = r
		}
		if greatest == i {
			// Greatest in correct position -> no need to swap further
			break
		}
		h.swap(i, greatest)
		i = greatest
	}
}

// Heapify converts an unordered slice into a valid max heap in-place
// Time: O(n)
func (h *MaxHeap) Heapify(items []int) {
	// I choose to copy data into this heap's items. ALternative is to take ownership of the array directly.
	h.items = make([]int, len(items))
	copy(h.items, items)

	// Hint: Start from last non-leaf node and sift down each node
	// Meaning pushing a node down the min heap until its smaller than both children
	// Last non-leaf is at index (n/2 - 1)
	for i := len(items)/2 - 1; 0 <= i; i-- {
		h.siftDown(i)
	}
}

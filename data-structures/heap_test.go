package datastructures

import "testing"

func TestMinHeap_PushPop(t *testing.T) {
	h := NewMinHeap()

	// Push elements
	h.Push(5)
	h.Push(3)
	h.Push(7)
	h.Push(1)
	h.Push(2)

	// Pop should return in ascending order
	expected := []int{1, 2, 3, 5, 7}
	for _, want := range expected {
		got, ok := h.Pop()
		if !ok {
			t.Fatalf("Pop returned not ok, expected %d", want)
		}
		if got != want {
			t.Errorf("Pop = %d, want %d", got, want)
		}
	}

	// Heap should be empty
	if !h.IsEmpty() {
		t.Error("Heap should be empty after popping all elements")
	}
}

func TestMinHeap_Peek(t *testing.T) {
	h := NewMinHeap()

	// Peek on empty heap
	_, ok := h.Peek()
	if ok {
		t.Error("Peek on empty heap should return not ok")
	}

	h.Push(5)
	h.Push(3)

	// Peek should return min without removing
	val, ok := h.Peek()
	if !ok {
		t.Error("Peek returned not ok")
	}
	if val != 3 {
		t.Errorf("Peek = %d, want 3", val)
	}
	if h.Size() != 2 {
		t.Errorf("Size after Peek = %d, want 2", h.Size())
	}
}

func TestMinHeap_PopEmpty(t *testing.T) {
	h := NewMinHeap()

	_, ok := h.Pop()
	if ok {
		t.Error("Pop on empty heap should return not ok")
	}
}

func TestMinHeap_Size(t *testing.T) {
	h := NewMinHeap()

	if h.Size() != 0 {
		t.Errorf("Size of new heap = %d, want 0", h.Size())
	}

	h.Push(1)
	h.Push(2)
	h.Push(3)

	if h.Size() != 3 {
		t.Errorf("Size after 3 pushes = %d, want 3", h.Size())
	}

	h.Pop()

	if h.Size() != 2 {
		t.Errorf("Size after pop = %d, want 2", h.Size())
	}
}

func TestMinHeap_Duplicates(t *testing.T) {
	h := NewMinHeap()

	h.Push(3)
	h.Push(3)
	h.Push(3)

	for i := 0; i < 3; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Fatalf("Pop %d returned not ok", i)
		}
		if val != 3 {
			t.Errorf("Pop %d = %d, want 3", i, val)
		}
	}
}

func TestMinHeap_SingleElement(t *testing.T) {
	h := NewMinHeap()

	h.Push(42)

	val, ok := h.Pop()
	if !ok || val != 42 {
		t.Errorf("Pop = (%d, %v), want (42, true)", val, ok)
	}

	if !h.IsEmpty() {
		t.Error("Heap should be empty")
	}
}

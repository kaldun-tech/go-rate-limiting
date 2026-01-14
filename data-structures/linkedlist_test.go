package datastructures

import (
	"reflect"
	"testing"
)

// Helper function to create a linked list from a slice
func createList(vals []int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	head := &ListNode{Val: vals[0]}
	curr := head
	for i := 1; i < len(vals); i++ {
		curr.Next = &ListNode{Val: vals[i]}
		curr = curr.Next
	}
	return head
}

// Helper function to convert linked list to slice for comparison
func listToSlice(head *ListNode) []int {
	result := []int{}
	for curr := head; curr != nil; curr = curr.Next {
		result = append(result, curr.Val)
	}
	return result
}

// TestReverseList tests the iterative reverse function
func TestReverseList(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "empty list",
			input: []int{},
			want:  []int{},
		},
		{
			name:  "single node",
			input: []int{1},
			want:  []int{1},
		},
		{
			name:  "two nodes",
			input: []int{1, 2},
			want:  []int{2, 1},
		},
		{
			name:  "multiple nodes",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createList(tt.input)
			reversed := ReverseList(head)
			got := listToSlice(reversed)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseList() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestReverseListRecursive tests the recursive reverse function
func TestReverseListRecursive(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "empty list",
			input: []int{},
			want:  []int{},
		},
		{
			name:  "single node",
			input: []int{1},
			want:  []int{1},
		},
		{
			name:  "two nodes",
			input: []int{1, 2},
			want:  []int{2, 1},
		},
		{
			name:  "multiple nodes",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createList(tt.input)
			reversed := ReverseListRecursive(head)
			got := listToSlice(reversed)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseListRecursive() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasCycle tests cycle detection
func TestHasCycle(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		head := createList([]int{1, 2, 3, 4, 5})
		if HasCycle(head) {
			t.Error("HasCycle() = true, want false for list without cycle")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		if HasCycle(nil) {
			t.Error("HasCycle() = true, want false for empty list")
		}
	})

	t.Run("single node no cycle", func(t *testing.T) {
		head := &ListNode{Val: 1}
		if HasCycle(head) {
			t.Error("HasCycle() = true, want false for single node")
		}
	})

	t.Run("cycle at end", func(t *testing.T) {
		head := createList([]int{1, 2, 3, 4})
		// Create cycle: 4 -> 2
		curr := head
		second := head.Next
		for curr.Next != nil {
			curr = curr.Next
		}
		curr.Next = second
		if !HasCycle(head) {
			t.Error("HasCycle() = false, want true for list with cycle")
		}
	})

	t.Run("cycle at beginning", func(t *testing.T) {
		head := createList([]int{1, 2, 3})
		// Create cycle: 3 -> 1
		curr := head
		for curr.Next != nil {
			curr = curr.Next
		}
		curr.Next = head
		if !HasCycle(head) {
			t.Error("HasCycle() = false, want true for list with cycle at beginning")
		}
	})

	t.Run("single node self cycle", func(t *testing.T) {
		head := &ListNode{Val: 1}
		head.Next = head
		if !HasCycle(head) {
			t.Error("HasCycle() = false, want true for self-cycle")
		}
	})
}

// TestFindMiddle tests finding the middle node
func TestFindMiddle(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		wantVal  int
		wantNil  bool
	}{
		{
			name:    "empty list",
			input:   []int{},
			wantNil: true,
		},
		{
			name:    "single node",
			input:   []int{1},
			wantVal: 1,
		},
		{
			name:    "two nodes - returns second",
			input:   []int{1, 2},
			wantVal: 2,
		},
		{
			name:    "odd length - middle node",
			input:   []int{1, 2, 3, 4, 5},
			wantVal: 3,
		},
		{
			name:    "even length - second middle",
			input:   []int{1, 2, 3, 4},
			wantVal: 3,
		},
		{
			name:    "six nodes",
			input:   []int{1, 2, 3, 4, 5, 6},
			wantVal: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createList(tt.input)
			middle := FindMiddle(head)
			if tt.wantNil {
				if middle != nil {
					t.Errorf("FindMiddle() = %v, want nil", middle)
				}
			} else {
				if middle == nil {
					t.Error("FindMiddle() = nil, want non-nil")
				} else if middle.Val != tt.wantVal {
					t.Errorf("FindMiddle() value = %d, want %d", middle.Val, tt.wantVal)
				}
			}
		})
	}
}

// TestMergeTwoLists tests merging two sorted lists
func TestMergeTwoLists(t *testing.T) {
	tests := []struct {
		name string
		l1   []int
		l2   []int
		want []int
	}{
		{
			name: "both empty",
			l1:   []int{},
			l2:   []int{},
			want: []int{},
		},
		{
			name: "first empty",
			l1:   []int{},
			l2:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "second empty",
			l1:   []int{1, 2, 3},
			l2:   []int{},
			want: []int{1, 2, 3},
		},
		{
			name: "same length",
			l1:   []int{1, 3, 5},
			l2:   []int{2, 4, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "different lengths",
			l1:   []int{1, 2, 4},
			l2:   []int{1, 3, 4, 5, 6},
			want: []int{1, 1, 2, 3, 4, 4, 5, 6},
		},
		{
			name: "first list shorter",
			l1:   []int{5},
			l2:   []int{1, 2, 4},
			want: []int{1, 2, 4, 5},
		},
		{
			name: "all first then second",
			l1:   []int{1, 2, 3},
			l2:   []int{4, 5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l1 := createList(tt.l1)
			l2 := createList(tt.l2)
			merged := MergeTwoLists(l1, l2)
			got := listToSlice(merged)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeTwoLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRemoveNthFromEnd tests removing the nth node from the end
func TestRemoveNthFromEnd(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
		want  []int
	}{
		{
			name:  "remove last node",
			input: []int{1, 2, 3, 4, 5},
			n:     1,
			want:  []int{1, 2, 3, 4},
		},
		{
			name:  "remove second from end",
			input: []int{1, 2, 3, 4, 5},
			n:     2,
			want:  []int{1, 2, 3, 5},
		},
		{
			name:  "remove head",
			input: []int{1, 2, 3, 4, 5},
			n:     5,
			want:  []int{2, 3, 4, 5},
		},
		{
			name:  "single node",
			input: []int{1},
			n:     1,
			want:  []int{},
		},
		{
			name:  "two nodes remove first",
			input: []int{1, 2},
			n:     2,
			want:  []int{2},
		},
		{
			name:  "two nodes remove second",
			input: []int{1, 2},
			n:     1,
			want:  []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createList(tt.input)
			result := RemoveNthFromEnd(head, tt.n)
			got := listToSlice(result)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveNthFromEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsPalindrome tests palindrome detection
func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  bool
	}{
		{
			name:  "empty list",
			input: []int{},
			want:  true,
		},
		{
			name:  "single node",
			input: []int{1},
			want:  true,
		},
		{
			name:  "two nodes palindrome",
			input: []int{1, 1},
			want:  true,
		},
		{
			name:  "two nodes not palindrome",
			input: []int{1, 2},
			want:  false,
		},
		{
			name:  "odd length palindrome",
			input: []int{1, 2, 3, 2, 1},
			want:  true,
		},
		{
			name:  "even length palindrome",
			input: []int{1, 2, 2, 1},
			want:  true,
		},
		{
			name:  "not palindrome",
			input: []int{1, 2, 3},
			want:  false,
		},
		{
			name:  "longer palindrome",
			input: []int{1, 2, 3, 4, 3, 2, 1},
			want:  true,
		},
		{
			name:  "almost palindrome",
			input: []int{1, 2, 3, 4, 2, 1},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createList(tt.input)
			got := IsPalindrome(head)
			if got != tt.want {
				t.Errorf("IsPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkReverseList(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Need to recreate list each time as it gets modified
		list := createList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		ReverseList(list)
	}
}

func BenchmarkReverseListRecursive(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := createList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		ReverseListRecursive(list)
	}
}

func BenchmarkHasCycle(b *testing.B) {
	head := createList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycle(head)
	}
}

func BenchmarkFindMiddle(b *testing.B) {
	head := createList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindMiddle(head)
	}
}

func BenchmarkMergeTwoLists(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l1 := createList([]int{1, 3, 5, 7, 9})
		l2 := createList([]int{2, 4, 6, 8, 10})
		MergeTwoLists(l1, l2)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list := createList([]int{1, 2, 3, 4, 5, 4, 3, 2, 1})
		IsPalindrome(list)
	}
}

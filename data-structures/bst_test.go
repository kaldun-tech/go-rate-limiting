package datastructures

import (
	"reflect"
	"testing"
)

func TestBST_InsertAndSearch(t *testing.T) {
	bst := NewBST()

	// Insert into empty tree
	if !bst.Insert(5) {
		t.Error("Insert(5) failed on empty tree")
	}

	// Search for existing value
	if !bst.Search(5) {
		t.Error("Search(5) should return true")
	}

	// Search for non-existing value
	if bst.Search(3) {
		t.Error("Search(3) should return false")
	}

	// Insert more values
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(4)

	if !bst.Search(1) || !bst.Search(3) || !bst.Search(4) || !bst.Search(7) {
		t.Error("Search failed for inserted values")
	}

	// Try to insert duplicate
	if bst.Insert(5) {
		t.Error("Insert(5) should return false for duplicate")
	}
}

func TestBST_InOrder(t *testing.T) {
	bst := NewBST()

	// Empty tree
	if got := bst.InOrder(); len(got) != 0 {
		t.Errorf("InOrder() on empty tree = %v, want []", got)
	}

	// Single node
	bst.Insert(5)
	if got := bst.InOrder(); !reflect.DeepEqual(got, []int{5}) {
		t.Errorf("InOrder() = %v, want [5]", got)
	}

	// Multiple nodes - should return sorted order
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(9)
	bst.Insert(4)

	want := []int{1, 3, 4, 5, 7, 9}
	if got := bst.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() = %v, want %v", got, want)
	}
}

func TestBST_Height(t *testing.T) {
	bst := NewBST()

	// Empty tree
	if got := bst.Height(); got != 0 {
		t.Errorf("Height() on empty tree = %d, want 0", got)
	}

	// Single node
	bst.Insert(5)
	if got := bst.Height(); got != 1 {
		t.Errorf("Height() = %d, want 1", got)
	}

	// Balanced tree: height 3
	//       5
	//      / \
	//     3   7
	//    /     \
	//   1       9
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(9)

	if got := bst.Height(); got != 3 {
		t.Errorf("Height() = %d, want 3", got)
	}

	// Unbalanced tree (right-skewed)
	bst2 := NewBST()
	bst2.Insert(1)
	bst2.Insert(2)
	bst2.Insert(3)
	bst2.Insert(4)

	if got := bst2.Height(); got != 4 {
		t.Errorf("Height() on unbalanced tree = %d, want 4", got)
	}
}

func TestBST_Delete_LeafNode(t *testing.T) {
	bst := NewBST()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	// Delete leaf node
	if !bst.Delete(3) {
		t.Error("Delete(3) should return true")
	}

	if bst.Search(3) {
		t.Error("Search(3) should return false after deletion")
	}

	// Verify tree structure is intact
	want := []int{5, 7}
	if got := bst.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() after delete = %v, want %v", got, want)
	}
}

func TestBST_Delete_OneChild(t *testing.T) {
	bst := NewBST()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)

	// Delete node with one child (3 has only left child 1)
	if !bst.Delete(3) {
		t.Error("Delete(3) should return true")
	}

	if bst.Search(3) {
		t.Error("Search(3) should return false after deletion")
	}

	// Child should still be present
	if !bst.Search(1) {
		t.Error("Search(1) should return true - child should remain")
	}

	want := []int{1, 5, 7}
	if got := bst.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() after delete = %v, want %v", got, want)
	}
}

func TestBST_Delete_TwoChildren(t *testing.T) {
	bst := NewBST()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)
	bst.Insert(1)
	bst.Insert(4)
	bst.Insert(6)
	bst.Insert(9)

	// Delete node with two children
	if !bst.Delete(5) {
		t.Error("Delete(5) should return true")
	}

	if bst.Search(5) {
		t.Error("Search(5) should return false after deletion")
	}

	// All other nodes should remain
	want := []int{1, 3, 4, 6, 7, 9}
	if got := bst.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() after delete = %v, want %v", got, want)
	}

	// Tree should still be valid
	if !bst.Search(1) || !bst.Search(3) || !bst.Search(4) {
		t.Error("Left subtree nodes should still exist")
	}
	if !bst.Search(6) || !bst.Search(7) || !bst.Search(9) {
		t.Error("Right subtree nodes should still exist")
	}
}

func TestBST_Delete_NonExistent(t *testing.T) {
	bst := NewBST()
	bst.Insert(5)
	bst.Insert(3)
	bst.Insert(7)

	// Try to delete non-existent value
	if bst.Delete(10) {
		t.Error("Delete(10) should return false for non-existent value")
	}

	// Tree should be unchanged
	want := []int{3, 5, 7}
	if got := bst.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() should be unchanged = %v, want %v", got, want)
	}
}

func TestBST_Delete_EmptyTree(t *testing.T) {
	bst := NewBST()

	// Delete from empty tree
	if bst.Delete(5) {
		t.Error("Delete(5) on empty tree should return false")
	}
}

func TestBST_Delete_Root(t *testing.T) {
	bst := NewBST()
	bst.Insert(5)

	// Delete the only node (root)
	if !bst.Delete(5) {
		t.Error("Delete(5) should return true")
	}

	if bst.Search(5) {
		t.Error("Search(5) should return false after deletion")
	}

	// Tree should be empty
	if got := bst.Height(); got != 0 {
		t.Errorf("Height() after deleting root = %d, want 0", got)
	}
}

func TestBST_ComplexOperations(t *testing.T) {
	bst := NewBST()

	// Build a tree
	values := []int{50, 30, 70, 20, 40, 60, 80, 10, 25, 35, 65}
	for _, v := range values {
		bst.Insert(v)
	}

	// Verify initial state
	if got := bst.Height(); got != 4 {
		t.Errorf("Initial height = %d, want 4", got)
	}

	// Delete a few nodes with different cases
	bst.Delete(20) // Node with two children
	bst.Delete(10) // Leaf node
	bst.Delete(30) // Node with two children

	// Verify remaining nodes
	remaining := []int{25, 35, 40, 50, 60, 65, 70, 80}
	if got := bst.InOrder(); !reflect.DeepEqual(got, remaining) {
		t.Errorf("InOrder() after deletions = %v, want %v", got, remaining)
	}

	// Verify all remaining values are searchable
	for _, v := range remaining {
		if !bst.Search(v) {
			t.Errorf("Search(%d) should return true", v)
		}
	}
}

// Benchmark Insert operations
func BenchmarkBST_Insert(b *testing.B) {
	bst := NewBST()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.Insert(i)
	}
}

// Benchmark Search operations
func BenchmarkBST_Search(b *testing.B) {
	bst := NewBST()
	for i := 0; i < 1000; i++ {
		bst.Insert(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.Search(i % 1000)
	}
}

// Benchmark Delete operations
func BenchmarkBST_Delete(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		bst := NewBST()
		for j := 0; j < 1000; j++ {
			bst.Insert(j)
		}
		b.StartTimer()
		bst.Delete(500)
		b.StopTimer()
	}
}

// Benchmark InOrder traversal
func BenchmarkBST_InOrder(b *testing.B) {
	bst := NewBST()
	for i := 0; i < 1000; i++ {
		bst.Insert(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.InOrder()
	}
}

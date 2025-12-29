package datastructures

// TreeNode represents a node in a binary search tree
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// BST implements a basic Binary Search Tree
// Common interview operations to implement:
// - Insert: Add a new value
// - Search: Find if a value exists
// - Delete: Remove a value (hardest - handle 3 cases)
// - InOrder: Return values in sorted order
// - Validate: Check if tree is a valid BST
type BST struct {
	root *TreeNode
}

// NewBST creates an empty binary search tree
func NewBST() *BST {
	return &BST{}
}

// Insert adds a value to the BST
// Returns true if inserted, false if value already exists
func (bst *BST) Insert(val int) bool {
	// TODO: Implement
	// Hint: Use recursive helper or iterative approach
	return false
}

// Search checks if a value exists in the BST
func (bst *BST) Search(val int) bool {
	// TODO: Implement
	return false
}

// Delete removes a value from the BST
// This is the trickiest operation - three cases:
// 1. Node has no children (leaf) - just remove
// 2. Node has one child - replace with child
// 3. Node has two children - replace with inorder successor/predecessor
func (bst *BST) Delete(val int) bool {
	// TODO: Implement
	return false
}

// InOrder returns values in sorted order (left -> root -> right)
func (bst *BST) InOrder() []int {
	// TODO: Implement
	return nil
}

// IsValid checks if the tree is a valid BST
// Common interview question: for each node, all left descendants < node < all right descendants
func (bst *BST) IsValid() bool {
	// TODO: Implement
	// Hint: Pass min/max bounds down the tree
	return false
}

// Height returns the height of the tree
func (bst *BST) Height() int {
	// TODO: Implement
	return 0
}

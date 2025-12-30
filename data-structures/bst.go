package datastructures

import (
	"math/rand"
)

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

func (bst *BST) doInsert(n *TreeNode, val int) bool {
	if n == nil {
		// Empty tree
		bst.root = &TreeNode{
			Val: val,
		}
		return true
	} else if val == n.Val {
		// Already contains this value
		return false
	} else if val < n.Val {
		// Search left
		if n.Left == nil {
			// Insert left
			n.Left = &TreeNode{
				Val: val,
			}
			return true
		} else {
			return bst.doInsert(n.Left, val)
		}
	} else {
		// Search right (n.Val < val)
		if n.Right == nil {
			// Insert right
			n.Right = &TreeNode{
				Val: val,
			}
			return true
		} else {
			return bst.doInsert(n.Right, val)
		}
	}
}

// Insert adds a value to the BST
// Returns true if inserted, false if value already exists
func (bst *BST) Insert(val int) bool {
	return bst.doInsert(bst.root, val)
}

func (bst *BST) doSearch(n *TreeNode, val int) bool {
	if n == nil {
		// Empty tree
		return false
	} else if n.Val == val {
		// Found
		return true
	} else if val < n.Val {
		// Search left
		return bst.doSearch(n.Left, val)
	} else if n.Val < val {
		// Search right
		return bst.doSearch(n.Right, val)
	}

	return false
}

// Search checks if a value exists in the BST
func (bst *BST) Search(val int) bool {
	return bst.doSearch(bst.root, val)
}

// One hop right then all the way left
func (bst *BST) successor(n *TreeNode) *TreeNode {
	if n == nil || n.Right == nil {
		return nil
	}
	suc := n.Right
	for suc.Left != nil {
		suc = suc.Left
	}
	return suc
}

// One hop left then all the way right
func (bst *BST) predecessor(n *TreeNode) *TreeNode {
	if n == nil || n.Left == nil {
		return nil
	}
	pred := n.Left
	for pred.Right != nil {
		pred = pred.Right
	}
	return pred
}

// Deletes by replacing with a predecessor or successor with equal probability
func (bst *BST) hibbardDelete(n *TreeNode) *TreeNode {
	// Generate a random int [0, 1] with 50% probability
	var random int = rand.Intn(2)
	if random == 0 {
		// Use successor if available
		if n.Right != nil {
			suc := bst.successor(n)
			n.Val = suc.Val
			n.Right = bst.deleteHelper(n.Right, suc.Val)
			return n
		} else if n.Left != nil {
			// Fall back to predecessor if no right child
			pred := bst.predecessor(n)
			n.Val = pred.Val
			n.Left = bst.deleteHelper(n.Left, pred.Val)
			return n
		} else {
			// Leaf node
			return nil
		}
	} else {
		// Prefer predecessor
		if n.Left != nil {
			pred := bst.predecessor(n)
			n.Val = pred.Val
			n.Left = bst.deleteHelper(n.Left, pred.Val)
			return n
		} else if n.Right != nil {
			// Use successor if no left child
			suc := bst.successor(n)
			n.Val = suc.Val
			n.Right = bst.deleteHelper(n.Right, suc.Val)
			return n
		} else {
			// Leaf node
			return nil
		}
	}
}

func (bst *BST) deleteHelper(n *TreeNode, val int) *TreeNode {
	if n == nil {
		// Not found
		return nil
	}
	if val < n.Val {
		// Search left
		n.Left = bst.deleteHelper(n.Left, val)
	} else if val == n.Val {
		// Found matching
		if n.Left == nil && n.Right == nil {
			// Just remove leaf node
			return nil
		} else if n.Left == nil {
			// No left child -> return right child
			return n.Right
		} else if n.Right == nil {
			// No right child -> return left child
			return n.Left
		} else {
			// Node with two children -> Hibbard deletion
			return bst.hibbardDelete(n)
		}

	} else {
		// Search right
		n.Right = bst.deleteHelper(n.Right, val)
	}
	// Return deleted node
	return n
}

// Delete removes a value from the BST with Hibbard deletion
// This is the trickiest operation - three cases:
// 1. Node has no children (leaf) - just remove
// 2. Node has one child - replace with child
// 3. Node has two children - replace with inorder successor/predecessor
func (bst *BST) Delete(val int) bool {
	if !bst.Search(val) {
		return false
	}
	bst.root = bst.deleteHelper(bst.root, val)
	return true
}

func (bst *BST) inOrderHelperRecursive(n *TreeNode, result *[]int) {
	if n == nil {
		// Empty node
		return
	}
	if n.Left != nil {
		// Go left
		bst.inOrderHelperRecursive(n.Left, result)
	}
	*result = append(*result, n.Val)
	if n.Right != nil {
		// Go right
		bst.inOrderHelperRecursive(n.Right, result)
	}
}

func (bst *BST) inOrderHelperStack(n *TreeNode) []int {
	var result []int
	stack := []*TreeNode{}
	current := bst.root

	// Iterate while we have a reference node or can pop one off the stack
	for current != nil || 0 < len(stack) {
		// Go as far left as possibl, pushing to the stack
		for current != nil {
			stack = append(stack, current)
			current = current.Left
		}
		// Pop off the stack and process
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, current.Val)
		// Go right
		current = current.Right
	}
	return result
}

func (bst *BST) inOrderHelper(n *TreeNode, result *[]int) {
	bst.inOrderHelperRecursive(n, result)
}

// InOrder returns values in sorted order (left -> root -> right)
func (bst *BST) InOrder() []int {
	result := []int{}
	bst.inOrderHelperRecursive(bst.root, &result)
	return result
}

// IsValid checks if the tree is a valid BST
// Common interview question: for each node, all left descendants < node < all right descendants
func (bst *BST) IsValid() bool {
	// Base case: Check whether empty
	if bst.root == nil {
		return true
	}

	// TODO: Implement
	// Hint: Pass min/max bounds down the tree
	return false
}

// Gets the height of a sub-tree with root n
// Why private BST helper method?
// 1. Idiomatic Go: Go favors simple structs with functions, not OOP-style methods on data
// 2. Consistency: Your other operations (doInsert, doSearch) already follow this pattern
// 3. Separation of concerns: TreeNode is just data; BST contains the operations
// 4. Flexibility: Easier to test and swap implementations
func (bst *BST) heightHelper(n *TreeNode) int {
	if n == nil {
		// Base case: no height
		return 0
	}

	return 1 + max(bst.heightHelper(n.Left), bst.heightHelper(n.Right))
}

// Height returns the height of the tree
func (bst *BST) Height() int {
	return bst.heightHelper(bst.root)
}

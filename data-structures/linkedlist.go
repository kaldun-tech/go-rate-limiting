package datastructures

// ListNode represents a node in a singly linked list
type ListNode struct {
	Val  int
	Next *ListNode
}

// Common linked list interview problems to implement:
// - Reverse a linked list (iterative and recursive)
// - Detect cycle (Floyd's algorithm - fast/slow pointers)
// - Find middle node (fast/slow pointers)
// - Merge two sorted lists
// - Remove Nth node from end
// - Check if palindrome

// ReverseList reverses a linked list and returns the new head
// Example: 1->2->3->4->5 becomes 5->4->3->2->1
func ReverseList(head *ListNode) *ListNode {
	// TODO: Implement iteratively
	// Hint: Use three pointers (prev, curr, next)
	return nil
}

// ReverseListRecursive reverses a linked list recursively
func ReverseListRecursive(head *ListNode) *ListNode {
	// TODO: Implement
	// Hint: Base case, then reverse rest and fix pointers
	return nil
}

// HasCycle detects if a linked list has a cycle
// Use Floyd's Cycle Detection (tortoise and hare)
func HasCycle(head *ListNode) bool {
	// TODO: Implement
	// Hint: Fast pointer moves 2x, slow pointer moves 1x
	return false
}

// FindMiddle returns the middle node of the list
// If two middle nodes, return the second one
func FindMiddle(head *ListNode) *ListNode {
	// TODO: Implement
	// Hint: Fast/slow pointers
	return nil
}

// MergeTwoLists merges two sorted linked lists
// Returns the head of the merged list
func MergeTwoLists(l1, l2 *ListNode) *ListNode {
	// TODO: Implement
	// Hint: Use a dummy node to simplify edge cases
	return nil
}

// RemoveNthFromEnd removes the nth node from the end
// Example: list=1->2->3->4->5, n=2 returns 1->2->3->5
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	// TODO: Implement
	// Hint: Two pointers with n gap between them
	return nil
}

// IsPalindrome checks if linked list values form a palindrome
// Example: 1->2->3->2->1 returns true
func IsPalindrome(head *ListNode) bool {
	// TODO: Implement
	// Hint: Find middle, reverse second half, compare
	return false
}

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
// Implement iteratively using three pointers (prev, curr, next)
func ReverseList(head *ListNode) *ListNode {
	// Prev is the new head which starts as null
	var prev, curr *ListNode

	// Save -> Reverse -> Advance
	for curr = head; curr != nil; {
		// Save link to next before breaking it
		next := curr.Next
		// Reverse the link
		curr.Next = prev
		// Move prev forward to current
		prev = curr
		// Move current forward
		curr = next
	}
	// Return new head
	return prev
}

// ReverseListRecursive reverses a linked list recursively
// Modify pointers in place rather than allocating new memory
func ReverseListRecursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		// Done - base case: empty list or single node
		return head
	}

	// Recurse immediately to go all the way to the end of the list
	newHead := ReverseListRecursive(head.Next)
	// Reverse the pointer - grandchild points to this
	head.Next.Next = head
	// Break old forward child link
	head.Next = nil
	// Bubble up the new head
	return newHead
}

// HasCycle detects if a linked list has a cycle
// Use Floyd's Cycle Detection (tortoise and hare)
func HasCycle(head *ListNode) bool {
	// Fast pointer moves 2x, slow pointer moves 1x
	fast, slow := head, head
	// Check that fast and its next are in bounds
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if fast == slow {
			// Both pointers in the same spot indicates cycle
			return true
		}
	}
	// No cycle found
	return false
}

// FindMiddle returns the middle node of the list
// If two middle nodes, return the second one
func FindMiddle(head *ListNode) *ListNode {
	// Fast pointer moves 2x, slow pointer moves 1x
	fast, slow := head, head
	// Check fast and its next are in bounds
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	// In the even case will return second middle
	return slow
}

// MergeTwoLists merges two sorted linked lists
// Returns the head of the merged list
// O(n + m) time for n = len(l1), m = len(l2)
func MergeTwoLists(l1, l2 *ListNode) *ListNode {
	// Use a dummy node to simplify edge cases
	var dummyHead *ListNode = &ListNode{}
	n := dummyHead

	for l1 != nil && l2 != nil {
		// Compare the two to get next val
		if l1.Val < l2.Val {
			// Use lesser l1
			n.Next = l1
			l1 = l1.Next
		} else {
			// Use lesser/equal l2
			n.Next = l2
			l2 = l2.Next
		}
		// Advance n
		n = n.Next
	}

	// Stop iterating and attach remaining list at the end to save time
	if l1 != nil {
		// Exhausted l2, add remainder of l1
		n.Next = l1
	} else if l2 != nil {
		// Exhausted l1, add remainder of l2
		n.Next = l2
	}

	// Return dummy head's next
	return dummyHead.Next
}

// RemoveNthFromEnd removes the nth node from the end
// Example: list=1->2->3->4->5, n=2 returns 1->2->3->5
// Use two pointers with n gap between them
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	// Use a dummy head to handle removing the existing head
	dummy := &ListNode{Next: head}
	first, second := dummy, dummy

	// Advance second by n+1 steps
	for range n + 1 {
		if second == nil {
			// n is greater than length of list
			return nil
		}
		second = second.Next
	}

	// Iterate BOTH pointers to find where second reaches the end
	for second != nil {
		first = first.Next
		second = second.Next
	}

	// First points to previous of the target node. Remove thetarget
	first.Next = first.Next.Next

	// Return updated head which may be empty
	return dummy.Next
}

// IsPalindrome checks if linked list values form a palindrome
// Example: 1->2->3->2->1 returns true, 1->2->2->1 returns true
// 1->2->3 returns flase
func IsPalindrome(head *ListNode) bool {
	// Find middle, reverse second half, compare
	middle := FindMiddle(head)
	reverseSecondHalf := ReverseList(middle)

	n := head
	m := reverseSecondHalf
	for n != nil && m != nil {
		if n.Val != m.Val {
			// Mismatch
			return false
		}
		// Advance both pointers
		n = n.Next
		m = m.Next
	}

	return true
}

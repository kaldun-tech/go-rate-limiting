# Linked List Interview Whiteboard Questions

Practice guide for demonstrating linked list understanding in technical interviews.

## Time Investment Guide

**Total Study Time:** ~6-8 hours to complete all exercises
**Recommended Pace:** 30-60 minutes daily = 8-12 days

**Quick Reference:**
- Core implementations (Q1-7): ~4.5 hours
- Practice exercises (1-6): ~2 hours
- Conceptual review: ~1 hour

**Daily 30-60 Minute Sessions:**
- Day 1: Iterative Reverse (45 min)
- Day 2: Recursive Reverse (45 min)
- Day 3: Cycle Detection + Find Middle (1 hour)
- Day 4: Merge Two Lists (45 min)
- Day 5: Remove Nth from End (30 min)
- Day 6: Palindrome Check (45 min)
- Day 7: Mixed practice + edge cases (1 hour)
- Day 8: Mock interview (1 hour)

## Core Implementation Questions

### 1. Reverse a Linked List (Iterative) (⏱️ 30-40 min)
**Question:** Reverse a singly linked list iteratively. Return the new head.

**Example:**
```
Input:  1 -> 2 -> 3 -> 4 -> 5 -> null
Output: 5 -> 4 -> 3 -> 2 -> 1 -> null
```

**Key Points to Cover:**
- Use three pointers: `prev`, `curr`, `next`
- Initialize `prev = nil` (new tail will point to nil)
- Iterate through list, reversing one pointer at a time
- Return `prev` at the end (new head)
- Time: O(n), Space: O(1)

**Detailed Algorithm:**
```
1. Initialize prev = nil, curr = head
2. While curr != nil:
   a. Save next: next = curr.Next
   b. Reverse pointer: curr.Next = prev
   c. Move prev forward: prev = curr
   d. Move curr forward: curr = next
3. Return prev (new head)
```

**Visual Walkthrough (1->2->3->null):**
```
Initial:
prev=nil, curr=1->2->3->null

Iteration 1:
  next = 2
  1.Next = nil       (reverse pointer)
  prev = 1, curr = 2
  Result: nil<-1  2->3->null

Iteration 2:
  next = 3
  2.Next = 1         (reverse pointer)
  prev = 2, curr = 3
  Result: nil<-1<-2  3->null

Iteration 3:
  next = nil
  3.Next = 2         (reverse pointer)
  prev = 3, curr = nil
  Result: nil<-1<-2<-3

Loop ends (curr == nil)
Return prev = 3
Final: 3->2->1->null
```

**Common Mistakes:**
- Forgetting to save `next` before reversing (you lose the rest of the list!)
- Returning `curr` instead of `prev` (curr is nil at the end)
- Not handling empty list or single node

**Follow-up:** Can you do this in O(1) space? (Yes, this is O(1)!)

---

### 2. Reverse a Linked List (Recursive) (⏱️ 35-45 min)
**Question:** Reverse a singly linked list using recursion. Return the new head.

**Key Points to Cover:**
- Base case: empty list or single node (return head)
- Recursive case: reverse rest of list, then fix current node
- Key insight: `head.Next.Next = head` reverses one link
- Must set `head.Next = nil` to avoid cycles
- Time: O(n), Space: O(n) due to call stack
- The new head is always the last node (bubbles up from base case)

**Detailed Algorithm:**
```go
func ReverseListRecursive(head *ListNode) *ListNode {
    // Base case: empty or single node
    if head == nil || head.Next == nil {
        return head
    }

    // Recurse to the end - newHead will be the last node
    newHead := ReverseListRecursive(head.Next)

    // Reverse the pointer between head and head.Next
    head.Next.Next = head  // Make next node point back
    head.Next = nil        // Break forward link to avoid cycle

    return newHead  // Bubble up the new head (last node)
}
```

**Visual Walkthrough (1->2->3->null):**

**Going down (recursion):**
```
Call 1: head=1
  Not base case → recurse with head.Next=2

  Call 2: head=2
    Not base case → recurse with head.Next=3

    Call 3: head=3
      head.Next = null (BASE CASE!)
      return 3
```

**Coming back up (unwinding):**
```
Back in Call 2 (head=2):
  newHead = 3
  head.Next.Next = head → 2.Next.Next = 2 → 3.Next = 2
  head.Next = nil → 2.Next = nil
  State: 3->2->null (1 still points to 2)
  return 3

Back in Call 1 (head=1):
  newHead = 3
  head.Next.Next = head → 1.Next.Next = 1 → 2.Next = 1
  head.Next = nil → 1.Next = nil
  State: 3->2->1->null
  return 3

Final result: 3->2->1->null
```

**The "Aha!" Moment:**
- When we're in a recursive call for node `2`, `head.Next` is node `3`
- After recursing, node `3` has been reversed with everything after it
- `head.Next.Next = head` means: "Hey node 3, point back to me (node 2)"
- `head.Next = nil` breaks the old forward link to prevent cycles
- We return the `newHead` from the deepest call (the original last node)

**Follow-ups:**
- What's the space complexity? (O(n) due to call stack)
- Why is iterative better? (O(1) space, no stack overflow risk)
- Can you trace through 1->2->null?

---

### 3. Detect Cycle (Floyd's Algorithm) (⏱️ 25-30 min)
**Question:** Detect if a linked list has a cycle using Floyd's Cycle Detection.

**Key Points to Cover:**
- Use two pointers: slow (moves 1 step) and fast (moves 2 steps)
- If there's a cycle, fast will eventually catch slow
- If no cycle, fast will reach nil
- Time: O(n), Space: O(1)

**Algorithm:**
```
1. Initialize slow = head, fast = head
2. While fast != nil and fast.Next != nil:
   a. slow = slow.Next
   b. fast = fast.Next.Next
   c. if slow == fast: return true (cycle detected)
3. Return false (fast reached end, no cycle)
```

**Why it works:**
- In a cycle, fast gains 1 node on slow per iteration
- Fast will eventually "lap" slow inside the cycle
- Like runners on a circular track

**Follow-up:** Find the start of the cycle (Floyd's algorithm part 2)

---

### 4. Find Middle Node (⏱️ 20-25 min)
**Question:** Find the middle node of a linked list. If two middle nodes exist, return the second one.

**Key Points to Cover:**
- Use fast/slow pointers (tortoise and hare)
- Slow moves 1 step, fast moves 2 steps
- When fast reaches end, slow is at middle
- Time: O(n), Space: O(1)

**Algorithm:**
```
1. Initialize slow = head, fast = head
2. While fast != nil and fast.Next != nil:
   a. slow = slow.Next
   b. fast = fast.Next.Next
3. Return slow (middle node)
```

**Examples:**
```
1->2->3->4->5->null
When fast reaches 5.Next (nil), slow is at 3

1->2->3->4->null
When fast reaches 4 (fast.Next=nil), slow is at 3
```

**Follow-up:** How to return the first middle for even-length lists?

---

### 5. Merge Two Sorted Lists (⏱️ 30-40 min)
**Question:** Merge two sorted linked lists into one sorted list.

**Key Points to Cover:**
- Use a dummy node to simplify head handling
- Compare values from both lists, attach smaller one
- Handle remaining nodes when one list is exhausted
- Time: O(n+m), Space: O(1)

**Algorithm:**
```
1. Create dummy node
2. Use tail pointer to build result
3. While both lists have nodes:
   a. Compare values
   b. Attach smaller node to tail
   c. Advance that list's pointer
   d. Advance tail
4. Attach remaining nodes from non-empty list
5. Return dummy.Next
```

**Follow-up:** Merge k sorted lists (use min heap)

---

### 6. Remove Nth Node From End (⏱️ 25-30 min)
**Question:** Remove the nth node from the end of the list in one pass.

**Key Points to Cover:**
- Use two pointers with n-gap between them
- Move fast pointer n steps ahead
- Move both until fast reaches end
- Slow is now at (n-1)th from end
- Use dummy node to handle removing head
- Time: O(n), Space: O(1)

**Algorithm:**
```
1. Create dummy node pointing to head
2. fast = dummy, slow = dummy
3. Move fast n+1 steps ahead
4. Move both until fast reaches end
5. slow.Next = slow.Next.Next (remove node)
6. Return dummy.Next
```

**Edge cases:**
- n equals list length (remove head)
- List has only one node
- n = 1 (remove tail)

**Follow-up:** What if n is larger than list length?

---

### 7. Check if Palindrome (⏱️ 35-45 min)
**Question:** Determine if a linked list's values form a palindrome.

**Key Points to Cover:**
- Find middle using fast/slow pointers
- Reverse second half
- Compare first half with reversed second half
- Optionally restore the list
- Time: O(n), Space: O(1)

**Algorithm:**
```
1. Find middle (fast/slow)
2. Reverse second half
3. Compare values from start and middle
4. (Optional) Restore list by reversing again
5. Return result
```

**Example:**
```
1->2->3->2->1 → palindrome
1->2->3->4->5 → not palindrome
```

**Follow-up:** Can you do it with O(n) space? (Use stack/array)

---

## Whiteboard Practice Exercises

### Exercise 1: Iterative Reversal - Step by Step
**Setup:** Draw 1->2->3->4->null

**Tasks:**
1. Initialize all pointers (prev, curr, next)
2. Manually trace through each iteration
3. Draw the state of all pointers after each iteration
4. Show final result

**What to practice:**
- Not losing track of the list when reversing
- Proper pointer updates in correct order
- Understanding when loop terminates
- What to return

---

### Exercise 2: Recursive Reversal - Call Stack
**Setup:** Draw 1->2->3->null

**Tasks:**
1. Draw the recursion tree (calls going down)
2. Mark the base case
3. Show the unwinding with pointer changes
4. Label what newHead is at each level
5. Explain why head.Next.Next = head works

**What to practice:**
- Visualizing recursive stack
- Understanding the unwinding phase
- Pointer manipulation during unwinding
- Why we return the same newHead up the stack

---

### Exercise 3: Reversal Comparison
**Question:** Code both iterative and recursive on whiteboard

**Discuss:**
- Time complexity of each (both O(n))
- Space complexity difference (O(1) vs O(n))
- Which is better for interviews?
- When might recursive be preferred?
- Risk of stack overflow for long lists

---

### Exercise 4: Reverse Sublist
**Challenge:** Reverse nodes from position m to n

**Example:**
```
Input: 1->2->3->4->5, m=2, n=4
Output: 1->4->3->2->5
```

**Hints:**
- Use iterative reversal in a specific range
- Need pointers to: before-m, node-m, node-n, after-n
- Reverse between m and n
- Reconnect the parts

---

### Exercise 5: Cycle Detection Trace
**Setup:** Draw a list with a cycle

**Tasks:**
1. Initialize slow and fast at head
2. Trace 5-6 iterations showing positions
3. Show when they meet
4. Prove fast will always catch slow in a cycle

**Extension:** Find the cycle start node

---

### Exercise 6: Palindrome Check Trace
**Setup:** Draw 1->2->3->2->1

**Tasks:**
1. Find middle using fast/slow
2. Reverse second half (draw it)
3. Compare both halves
4. Show restoration (optional)

---

## Common Pitfalls

### Mistake 1: Losing the List
```go
// WRONG - loses reference to rest of list
func reverse(head *ListNode) *ListNode {
    prev := nil
    for head != nil {
        head.Next = prev  // LOST! Can't advance to next
        prev = head
        head = ???  // What goes here?
    }
}

// RIGHT - save next before modifying
func reverse(head *ListNode) *ListNode {
    prev := nil
    for head != nil {
        next := head.Next  // Save it first!
        head.Next = prev
        prev = head
        head = next
    }
    return prev
}
```

### Mistake 2: Returning Wrong Pointer
```go
// WRONG - curr is nil at loop end
func reverse(head *ListNode) *ListNode {
    var prev *ListNode
    for curr := head; curr != nil; {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    return curr  // curr is nil here!
}

// RIGHT - prev is the new head
func reverse(head *ListNode) *ListNode {
    var prev *ListNode
    for curr := head; curr != nil; {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    return prev  // prev is the new head
}
```

### Mistake 3: Recursive Without Breaking Link
```go
// WRONG - creates a cycle!
func reverse(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }
    newHead := reverse(head.Next)
    head.Next.Next = head
    // Missing: head.Next = nil
    return newHead
}

// Result: ... <- 2 <-> 1 (cycle!)
```

### Mistake 4: Fast/Slow Pointer Nil Check
```go
// WRONG - can panic
func findMiddle(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil {
        fast = fast.Next.Next  // Panic if fast.Next is nil!
        slow = slow.Next
    }
    return slow
}

// RIGHT - check both
func findMiddle(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        fast = fast.Next.Next
        slow = slow.Next
    }
    return slow
}
```

---

## Conceptual Questions

### Q1: Iterative vs Recursive Reversal
**Question:** When would you prefer iterative over recursive?

**Answer:**
- **Iterative preferred when:**
  - Large lists (avoid stack overflow)
  - Production code (O(1) space)
  - Performance critical (no function call overhead)

- **Recursive preferred when:**
  - Teaching/learning (clearer logic)
  - Short lists (stack safe)
  - Already recursive problem (DFS on tree)

---

### Q2: Why Three Pointers?
**Question:** In iterative reversal, why do we need prev, curr, AND next?

**Answer:**
- **prev:** The new "next" for current node (building reversed list)
- **curr:** The node we're currently processing
- **next:** Saved reference to continue forward (before we break curr.Next)

Without all three, we either:
- Lose the rest of the list, OR
- Can't reverse the current link, OR
- Can't advance forward

---

### Q3: Dummy Node Pattern
**Question:** Why use a dummy node in linked list problems?

**Answer:**
- Simplifies edge cases (empty list, single node)
- Avoids special handling for head changes
- Provides stable pointer to return
- Used in: merge, remove, insertion problems

**Example:**
```go
dummy := &ListNode{}
tail := dummy
// Build list by appending to tail
return dummy.Next  // Skip dummy
```

---

## Complexity Cheat Sheet

| Operation | Time | Space | Notes |
|-----------|------|-------|-------|
| Reverse (iterative) | O(n) | O(1) | Preferred |
| Reverse (recursive) | O(n) | O(n) | Call stack |
| Detect cycle | O(n) | O(1) | Floyd's algorithm |
| Find middle | O(n) | O(1) | Fast/slow pointers |
| Merge two sorted | O(n+m) | O(1) | In-place |
| Remove nth from end | O(n) | O(1) | Two pointers |
| Check palindrome | O(n) | O(1) | Reverse half |

---

## Interview Tips

1. **Draw it out:** Always sketch the list before coding
2. **Edge cases:** Empty, single node, two nodes
3. **Pointer order:** Save before modifying (next = curr.Next comes first!)
4. **Nil checks:** Check before accessing .Next or .Val
5. **Dummy nodes:** Use for problems that modify head
6. **Fast/slow:** Remember to check both fast != nil AND fast.Next != nil
7. **Recursion trade-off:** Mention O(n) space vs O(1) for iterative
8. **Test your code:** Walk through with 1->2->null at minimum

---

## Practice Drill Schedule

**Day 1-2: Iterative Reversal**
- 5 times on whiteboard with 3-4 node lists
- Trace every pointer update
- Practice with nil, single node, two nodes

**Day 3-4: Recursive Reversal**
- Draw call stack for 3-4 node lists
- Explain unwinding phase out loud
- Compare with iterative approach

**Day 5: Fast/Slow Patterns**
- Cycle detection
- Find middle
- Remove nth from end

**Day 6: Integration**
- Palindrome check (combines reverse + fast/slow)
- Merge sorted lists
- Reverse sublist

**Day 7: Mock Interview**
- Random problem from above
- 30 minutes timed
- Explain, code, test

---

## Key Takeaways

1. **Master the three-pointer dance** for iterative reversal
2. **Understand recursion unwinding** for recursive reversal
3. **Fast/slow pointers** solve many list problems elegantly
4. **Always save next** before reversing pointers
5. **Draw first, code second** - visualize before implementing
6. **Edge cases matter** - test with nil, 1 node, 2 nodes
7. **Space complexity** - iterative is O(1), recursive is O(n)
8. **Practice both** - interviews might specify approach

Good luck!

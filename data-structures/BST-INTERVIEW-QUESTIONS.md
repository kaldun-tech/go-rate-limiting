# BST Interview Whiteboard Questions

Practice guide for demonstrating BST understanding in technical interviews.

## Time Investment Guide

**Total Study Time:** ~8-10 hours to complete all exercises
**Recommended Pace:** 30-60 minutes daily = 10-14 days

**Quick Reference:**
- Core implementations (Q1-6): ~5 hours (implement + practice)
- Conceptual questions (Q7-10): ~1 hour (review + understand)
- Coding challenges (Q11-15): ~3 hours (solve + optimize)
- Practice drills: ~1-2 hours (timed practice)

**Daily 30-60 Minute Sessions:**
- Day 1: Insert + Search (1 hour)
- Day 2: Delete (1 hour - most complex!)
- Day 3: InOrder + Height (45 min)
- Day 4: IsValid (1 hour - tricky!)
- Day 5: Conceptual review (30 min)
- Day 6-7: Coding challenges (1 hour each)
- Day 8: Mock interview practice (1 hour)

## Core Implementation Questions

### 1. Insert (⏱️ 20-30 min)
**Question:** Implement insert for a BST. How do you handle duplicates?

**Key Points to Cover:**
- Recursive vs iterative approach
- Base case: empty tree (create root)
- Compare and navigate left/right
- Duplicate handling (return false, ignore, or allow)
- Time: O(h) where h is height, Space: O(h) for recursion

**Follow-up:** What if we want to maintain a count for duplicate values?

---

### 2. Search (⏱️ 15-20 min)
**Question:** Implement search to find if a value exists in the BST.

**Key Points to Cover:**
- Similar navigation to insert
- Base cases: null (not found), value match (found)
- Binary search property enables pruning half the tree each step
- Time: O(h), Space: O(h) recursive or O(1) iterative

**Follow-up:** How would you find the kth smallest element?

---

### 3. Delete (Hibbard Deletion) (⏱️ 45-60 min)
**Question:** Implement delete for a BST. This is the trickiest operation.

**Key Points to Cover:**
- Three cases:
  1. Leaf node: just remove (return nil)
  2. One child: replace node with child
  3. Two children: replace with inorder successor OR predecessor
- Finding successor: go right once, then left as far as possible
- Finding predecessor: go left once, then right as far as possible
- Why randomize? Prevents tree from becoming unbalanced over time
- Time: O(h), Space: O(h)

**Follow-ups:**
- Why do we need special handling for two children?
- What's the difference between successor and predecessor?
- What happens to tree balance with repeated deletions?

---

### 4. InOrder Traversal (⏱️ 20-25 min)
**Question:** Return all values in sorted order.

**Key Points to Cover:**
- InOrder = Left → Root → Right
- Produces sorted output for BST (key property!)
- Recursive: simple and clean
- Iterative: requires explicit stack
- Time: O(n), Space: O(n) for result array + O(h) for stack

**Follow-ups:**
- Implement iteratively using a stack
- How would you do PreOrder? PostOrder?
- Can you do it with O(1) extra space? (Morris traversal)

---

### 5. Validate BST (⏱️ 30-40 min)
**Question:** Check if a binary tree is a valid BST.

**Key Points to Cover:**
- Common mistake: only checking immediate children
- Correct approach: track min/max bounds through recursion
- Each node must satisfy: min < node.val < max
- Left subtree gets updated max bound
- Right subtree gets updated min bound
- Time: O(n), Space: O(h)

**Alternative approach:**
- InOrder traversal should produce strictly increasing values
- Time: O(n), Space: O(n)

**Follow-up:** Which approach is better and why?

---

### 6. Height/Depth (⏱️ 15-20 min)
**Question:** Find the height of the BST.

**Key Points to Cover:**
- Height = longest path from root to leaf
- Recursive: 1 + max(height(left), height(right))
- Base case: null node has height 0
- Time: O(n), Space: O(h)

**Follow-ups:**
- What's the difference between height and depth?
- Height of single node? (1)
- Height of empty tree? (0)

---

## Conceptual Questions

### 7. Time Complexity Analysis
**Question:** What are the time complexities for BST operations?

**Answer:**
- **Best/Average case:** O(log n) - balanced tree
- **Worst case:** O(n) - degenerate tree (linked list)
- Operations: Insert, Delete, Search all have same complexity

**Follow-up:** How do balanced trees (AVL, Red-Black) solve the worst case?

---

### 8. BST Properties
**Question:** What makes a tree a valid BST?

**Answer:**
- For every node N:
  - All values in left subtree < N.val
  - All values in right subtree > N.val
- This must hold for ALL nodes recursively

**Follow-up:** Is this a valid BST?
```
    5
   / \
  3   7
 / \
1   6
```
**No!** 6 is in left subtree of 5, but 6 > 5.

---

### 9. When to Use BST
**Question:** When would you choose a BST over other data structures?

**Answer:**
- Need sorted order frequently
- Range queries (find all values between x and y)
- Dynamic dataset with frequent insertions/deletions
- Floor/ceiling operations
- Not ideal if: purely random access, need guaranteed O(log n), hash-based lookup sufficient

---

### 10. Successor/Predecessor
**Question:** Find the inorder successor of a given node.

**Key Points to Cover:**
- If node has right child: successor is leftmost node in right subtree
- If no right child: successor is lowest ancestor where node is in left subtree
- Without parent pointers: harder, need to track path during traversal
- Time: O(h)

---

## Coding Challenges

### 11. Lowest Common Ancestor
**Question:** Find the LCA of two nodes in a BST.

**Hint:** Use BST property - if both values are less than current, go left; if both greater, go right; otherwise current node is LCA.

**Time:** O(h)

---

### 12. Range Sum Query
**Question:** Find sum of all values between low and high (inclusive).

**Hint:** Use BST property to prune branches. If node.val < low, skip left subtree. If node.val > high, skip right subtree.

**Time:** O(n) worst case, but prunes effectively

---

### 13. Convert Sorted Array to BST
**Question:** Given a sorted array, create a balanced BST.

**Hint:** Use middle element as root recursively.

**Time:** O(n)

---

### 14. Find Mode
**Question:** Find the most frequently occurring value(s) in a BST.

**Hint:** InOrder traversal with tracking of current value count.

**Time:** O(n)

---

### 15. Serialize/Deserialize
**Question:** Convert BST to string and back.

**Approaches:**
- PreOrder traversal (can reconstruct BST without null markers)
- Level order with null markers
- InOrder + PreOrder (can uniquely reconstruct)

---

## Common Pitfalls

### Mistake 1: Forgetting Nil Checks
```go
// WRONG
func search(n *TreeNode, val int) bool {
    if n.Val == val {  // Panic if n is nil!
        return true
    }
    ...
}

// RIGHT
func search(n *TreeNode, val int) bool {
    if n == nil {
        return false
    }
    if n.Val == val {
        return true
    }
    ...
}
```

### Mistake 2: Invalid BST Validation
```go
// WRONG - only checks immediate children
func isValid(n *TreeNode) bool {
    if n.Left != nil && n.Left.Val >= n.Val {
        return false
    }
    if n.Right != nil && n.Right.Val <= n.Val {
        return false
    }
    return isValid(n.Left) && isValid(n.Right)
}

// This would incorrectly validate:
//     5
//    / \
//   3   7
//  / \
// 1   6  <- 6 > 5, invalid!
```

### Mistake 3: Height vs Size Confusion
```go
// Height: longest path to leaf
func height(n *TreeNode) int {
    if n == nil {
        return 0
    }
    return 1 + max(height(n.Left), height(n.Right))
}

// Size: total number of nodes
func size(n *TreeNode) int {
    if n == nil {
        return 0
    }
    return 1 + size(n.Left) + size(n.Right)  // Note: + not max!
}
```

---

## Design Discussion Questions

### 16. BST vs Hash Table
**When would you choose each?**

**BST:**
- Need sorted order
- Range queries
- Floor/ceiling operations
- Predictable O(log n) with balancing

**Hash Table:**
- Only need insert/search/delete
- Don't care about order
- Want average O(1) operations
- Have good hash function

---

### 17. Balancing Strategy
**Question:** Your BST is getting unbalanced. What are your options?

**Answers:**
- Rebuild periodically from sorted traversal
- Switch to self-balancing tree (AVL, Red-Black)
- Use randomization (treap)
- Consider B-tree for disk storage

---

## Practice Drill

**Warm-up (5 minutes each):**
1. Draw and explain Insert(10), Insert(5), Insert(15), Insert(3)
2. Delete node with two children - draw before/after
3. Validate this tree: Is it a BST? Why/why not?

**Medium (10-15 minutes each):**
4. Implement IsValid with bounds checking on whiteboard
5. Implement Delete with all three cases
6. Find kth smallest element

**Advanced (20+ minutes):**
7. Serialize/Deserialize a BST
8. Find all paths that sum to target value
9. Convert BST to sorted doubly linked list (in-place)

---

## Key Takeaways for Interviews

1. **Always check for nil** before accessing node properties
2. **BST property** - all left descendants < node < all right descendants
3. **Recursion is your friend** - most BST operations are naturally recursive
4. **Know the trade-offs** - BST vs other data structures
5. **Communicate clearly** - explain your approach before coding
6. **Test edge cases** - empty tree, single node, all left/right children
7. **Complexity analysis** - always be ready to discuss time/space complexity
8. **Ask clarifying questions** - duplicates allowed? balanced tree? constraints?

Good luck!

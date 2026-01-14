# Trie (Prefix Tree) Interview Whiteboard Questions

Practice guide for demonstrating Trie understanding in technical interviews.

## Time Investment Guide

**Total Study Time:** ~9-11 hours to complete all exercises
**Recommended Pace:** 30-60 minutes daily = 12-16 days

**Quick Reference:**
- Core implementations (Q1-6): ~5 hours
- Conceptual questions (Q7-10): ~1.5 hours
- Coding challenges (Q11-18): ~4 hours
- Practice drills: ~1.5 hours

**Daily 30-60 Minute Sessions:**
- Day 1: Trie structure + Insert (1 hour)
- Day 2: Search + StartsWith (45 min)
- Day 3: Delete (1 hour - complex!)
- Day 4: FindAllWithPrefix/Autocomplete (1 hour)
- Day 5: Longest Common Prefix (30 min)
- Day 6-7: Conceptual questions (30-45 min each)
- Day 8-11: Coding challenges (1 hour each)
- Day 12: Mock interview (1 hour)

## Core Implementation Questions

### 1. Insert (⏱️ 30-40 min)
**Question:** Implement insert to add a word to the trie.

**Key Points to Cover:**
- Start at root (dummy node)
- For each character in word:
  - If child doesn't exist, create new node
  - Move to child node
- Mark final node with `isEnd = true`
- Time: O(m) where m is word length
- Space: O(m) worst case (new nodes)

**Implementation approach:**
```go
func (t *Trie) Insert(word string) {
    node := t.root
    for _, ch := range word {  // Use rune for Unicode
        if node.children[ch] == nil {
            node.children[ch] = &TrieNode{
                children: make(map[rune]*TrieNode),
            }
        }
        node = node.children[ch]
    }
    node.isEnd = true
}
```

**Follow-up:** How would you modify this to count word frequency?

---

### 2. Search (⏱️ 20-25 min)
**Question:** Implement search to check if a complete word exists.

**Key Points to Cover:**
- Traverse character by character
- If any character not found, return false
- Must check `isEnd` at final node (distinguishes "car" from "card")
- Time: O(m), Space: O(1)

**Common mistake:** Forgetting to check `isEnd`

**Example:**
- Trie contains: "car", "card"
- `Search("car")` → true (car is marked as end)
- `Search("ca")` → false (ca exists but not marked as end)

**Follow-up:** How would you implement wildcard search with "." matching any character?

---

### 3. StartsWith (Prefix Search) (⏱️ 15-20 min)
**Question:** Check if any word starts with given prefix.

**Key Points to Cover:**
- Similar to Search but DON'T check `isEnd`
- Just verify path exists
- Time: O(m), Space: O(1)

**Key difference from Search:**
```go
// Search: must be complete word
func Search(word) {
    node := traverse(word)
    return node != nil && node.isEnd  // Must be end!
}

// StartsWith: just needs path
func StartsWith(prefix) {
    node := traverse(prefix)
    return node != nil  // Don't care about isEnd
}
```

**Follow-up:** How would you return count of words with this prefix?

---

### 4. Delete (⏱️ 45-60 min)
**Question:** Remove a word from the trie. Clean up unused nodes.

**Key Points to Cover:**
- Recursive approach is cleanest
- Base case: word not found, return false
- Mark `isEnd = false` at word end
- **Cleanup:** Only delete node if:
  - No children (not a prefix for other words)
  - Not end of another word
- Time: O(m), Space: O(m) for recursion

**Three scenarios:**
1. Word is prefix of longer word: Just unmark `isEnd`
2. Word has no children: Delete up to divergence point
3. Word shares prefix: Delete only unique suffix

**Example:**
```
Trie: "car", "card", "cat"
Delete("card"):
  - Unmark 'd' as end
  - Delete 'd' (no children, not end)
  - Keep "car" nodes (c-a-r is end of another word)
```

**Follow-up:** Can you implement iteratively instead of recursively?

---

### 5. FindAllWithPrefix (Autocomplete) (⏱️ 45-60 min)
**Question:** Return all words starting with given prefix.

**Key Points to Cover:**
- Two phases:
  1. Navigate to prefix node
  2. DFS/BFS to collect all words from there
- Need to build words during traversal
- Time: O(p + n*k) where p=prefix len, n=results, k=avg word len
- Space: O(n*k) for results + O(k) for recursion

**Implementation approach:**
```go
func (t *Trie) FindAllWithPrefix(prefix string) []string {
    // Phase 1: Find prefix node
    node := t.root
    for _, ch := range prefix {
        if node.children[ch] == nil {
            return []string{}  // Prefix doesn't exist
        }
        node = node.children[ch]
    }

    // Phase 2: DFS from prefix node
    results := []string{}
    t.dfs(node, prefix, &results)
    return results
}

func (t *Trie) dfs(node *TrieNode, path string, results *[]string) {
    if node.isEnd {
        *results = append(*results, path)
    }
    for ch, child := range node.children {
        t.dfs(child, path + string(ch), results)
    }
}
```

**Follow-up:** How would you return only top K most frequent words?

---

### 6. Longest Common Prefix (⏱️ 25-30 min)
**Question:** Find the longest common prefix of all words in the trie.

**Key Points to Cover:**
- Traverse from root while:
  - Only one child exists
  - Not at word end (unless continuing)
- Stop when branch or end
- Time: O(m) where m is result length
- Space: O(1)

**Follow-up:** What if trie is empty?

---

## Conceptual Questions

### 7. Time Complexity Analysis
**Question:** What are the time complexities for trie operations?

**Answer:**
- **Insert:** O(m) - m is word length
- **Search:** O(m)
- **StartsWith:** O(m)
- **Delete:** O(m)
- **FindAllWithPrefix:** O(p + n*k) - p=prefix, n=results, k=avg length

**Space:**
- Total: O(ALPHABET_SIZE * N * M) worst case
- In practice: Much better due to shared prefixes
- English lowercase: 26 pointers per node max

**Follow-up:** How does this compare to binary search on sorted array?

---

### 8. Trie Properties
**Question:** What makes a trie efficient for string operations?

**Answer:**
- **Prefix sharing** - Common prefixes stored once
- **No collisions** - Unlike hash tables
- **Sorted output** - In-order traversal gives sorted words
- **Fast prefix queries** - O(prefix length) not O(n)
- **Predictable** - No worst case like unbalanced BST

**Trade-offs:**
- **More space** - Pointer overhead per node
- **Cache unfriendly** - Pointer chasing
- **Alphabet dependent** - Large alphabets = more space

**Follow-up:** When would a hash set be better than a trie?

---

### 9. When to Use Trie
**Question:** When would you choose a trie over other data structures?

**Answer:**

**Use Trie when:**
- Autocomplete/typeahead
- Spell checking
- IP routing (longest prefix match)
- Phone directory (prefix search)
- Word games (Boggle, Scrabble validation)
- DNA sequence analysis

**Use Hash Set when:**
- Only exact match needed
- No prefix queries
- Memory constrained
- Random strings (no common prefixes)

**Use Binary Search when:**
- Read-only or infrequent updates
- Need range queries
- Memory is critical

---

### 10. Alphabet Size Optimization
**Question:** How would you optimize a trie for different alphabets?

**Answer:**

**Small, fixed alphabet (a-z):**
```go
type TrieNode struct {
    children [26]*TrieNode  // Array, not map
    isEnd    bool
}
// Access: children[ch - 'a']
```

**Large/Unicode alphabet:**
```go
type TrieNode struct {
    children map[rune]*TrieNode  // Map for sparse storage
    isEnd    bool
}
```

**Binary alphabet (0/1):**
```go
type TrieNode struct {
    zero, one *TrieNode  // Just two pointers
    isEnd     bool
}
```

**Follow-up:** What about case-insensitive search?

---

## Coding Challenges

### 11. Word Search II (Boggle)
**Question:** Given a 2D board and word list, find all words on the board.

**Hint:**
- Build trie of dictionary words
- DFS from each cell, pruning with trie
- Much faster than searching each word individually

**Time:** O(M*N*4^L) where L is max word length

**Key optimization:** Trie allows early pruning when prefix doesn't exist!

---

### 12. Replace Words
**Question:** Given dictionary of roots and sentence, replace words with roots.

**Example:**
- Roots: ["cat", "bat", "rat"]
- Sentence: "the cattle was rattled by the battery"
- Output: "the cat was rat by the bat"

**Hint:**
- Build trie of roots
- For each word, find shortest prefix in trie
- Time: O(N) where N is total characters

---

### 13. Palindrome Pairs
**Question:** Find all pairs of words that form palindromes when concatenated.

**Example:**
- Words: ["bat", "tab", "cat"]
- Output: [[0,1], [1,0]] (bat+tab, tab+bat)

**Hint:** Use trie with reverse insertion. Advanced problem!

---

### 14. Design Add and Search Words Data Structure
**Question:** Support adding words and searching with '.' wildcard.

**Key Points:**
- Insert: Same as normal trie
- Search with '.': Try all children branches
- Time: O(m) for no wildcards, O(26^w * m) with w wildcards

**Implementation:**
```go
func search(node *TrieNode, word string, idx int) bool {
    if idx == len(word) {
        return node.isEnd
    }

    if word[idx] == '.' {
        // Try all children
        for _, child := range node.children {
            if search(child, word, idx+1) {
                return true
            }
        }
        return false
    }

    // Normal character
    child := node.children[rune(word[idx])]
    if child == nil {
        return false
    }
    return search(child, word, idx+1)
}
```

---

### 15. Maximum XOR of Two Numbers
**Question:** Find maximum XOR of any two numbers in array.

**Hint:** Build binary trie (bits as path). For each number, greedily pick opposite bit when possible.

**Time:** O(32*N) for 32-bit integers

**Key insight:** Trie works on bits, not just strings!

---

### 16. Word Squares
**Question:** Find all word squares (n x n grid where rows = columns).

**Hint:** Backtrack with trie pruning - check if column prefixes exist.

---

### 17. Top K Frequent Words
**Question:** Find k most frequent words with autocomplete.

**Hint:**
- Extend TrieNode with frequency counter
- Use min-heap of size K during DFS
- Return in frequency order (break ties alphabetically)

---

### 18. Longest Word in Dictionary
**Question:** Find longest word that can be built one character at a time.

**Hint:** DFS in trie, only traverse if intermediate nodes are word ends.

---

## Common Pitfalls

### Mistake 1: Forgetting isEnd Flag
```go
// WRONG - can't distinguish prefix from word
type TrieNode struct {
    children map[rune]*TrieNode
}

// Trie: "car", "card"
// Search("car") vs StartsWith("car") - no difference!

// RIGHT
type TrieNode struct {
    children map[rune]*TrieNode
    isEnd    bool  // Essential!
}
```

---

### Mistake 2: Not Initializing Children Map
```go
// WRONG
func Insert(word string) {
    node := t.root
    for _, ch := range word {
        node.children[ch] = &TrieNode{}  // children is nil!
    }
}

// RIGHT
func Insert(word string) {
    node := t.root
    for _, ch := range word {
        if node.children[ch] == nil {
            node.children[ch] = &TrieNode{
                children: make(map[rune]*TrieNode),  // Initialize!
            }
        }
        node = node.children[ch]
    }
}
```

---

### Mistake 3: Using Byte Instead of Rune for Unicode
```go
// WRONG - breaks on Unicode
for i := 0; i < len(word); i++ {
    ch := word[i]  // byte - only one UTF-8 byte!
}

// RIGHT - handles multi-byte characters
for _, ch := range word {
    // ch is rune (int32)
}
```

**Example:** "hello世界"
- Byte iteration: gets individual UTF-8 bytes (broken Chinese)
- Rune iteration: gets complete Unicode code points

---

### Mistake 4: Incomplete Delete
```go
// WRONG - leaves unused nodes
func Delete(word string) {
    node := find(word)
    if node != nil {
        node.isEnd = false  // Doesn't clean up!
    }
}

// RIGHT - remove childless nodes
func delete(node *TrieNode, word string, idx int) bool {
    if idx == len(word) {
        if !node.isEnd {
            return false  // Word doesn't exist
        }
        node.isEnd = false
        return len(node.children) == 0  // Delete if no children
    }

    ch := rune(word[idx])
    child := node.children[ch]
    if child == nil {
        return false
    }

    shouldDelete := delete(child, word, idx+1)

    if shouldDelete {
        delete(node.children, ch)
        return !node.isEnd && len(node.children) == 0
    }
    return false
}
```

---

### Mistake 5: Building String Inefficiently
```go
// WRONG - creates new string each time (O(n^2))
func dfs(node *TrieNode, path string) {
    for ch, child := range node.children {
        dfs(child, path + string(ch))  // String concatenation!
    }
}

// BETTER - use strings.Builder
func dfs(node *TrieNode, builder *strings.Builder) {
    for ch, child := range node.children {
        builder.WriteRune(ch)
        dfs(child, builder)
        builder.Reset()  // Backtrack
    }
}
```

---

## Design Discussion Questions

### 19. Compressed Trie (Radix Tree)
**Question:** How would you reduce space usage?

**Answer:**
- Compress chains of single children
- Store edge labels (strings) instead of single characters
- Example: "test", "testing" shares "test" edge

**Trade-offs:**
- Less space (fewer nodes)
- More complex implementation
- Still O(m) operations

**Use cases:** Git objects, routing tables, Patricia trie

---

### 20. Trie vs Hash Table
**Question:** Compare tries to hash tables for string operations.

**Answer:**

| Operation | Trie | Hash Table |
|-----------|------|------------|
| Exact search | O(m) | O(1) avg |
| Prefix search | O(m) | O(n) |
| Insert | O(m) | O(1) avg |
| Space | O(ALPHABET * N * M) | O(N * M) |
| Worst case | O(m) | O(n) |
| Sorted iteration | Yes | No |

**Trie wins:** Prefix queries, autocomplete, sorted order
**Hash wins:** Exact lookup, space efficiency, simplicity

---

### 21. Case Sensitivity
**Question:** How would you handle case-insensitive search?

**Options:**

1. **Normalize on insert:**
```go
func Insert(word string) {
    t.insert(strings.ToLower(word))
}
```

2. **Double storage:**
```go
type TrieNode struct {
    children map[rune]*TrieNode  // Both 'A' and 'a'
}
```

3. **Ignore case on comparison:**
```go
// Convert 'A'-'Z' to 'a'-'z' during traversal
ch = unicode.ToLower(ch)
```

**Recommendation:** Option 1 (normalize) - simplest and space efficient

---

### 22. Memory Optimization
**Question:** Your trie is using too much memory. How to optimize?

**Options:**

1. **Array for small alphabets:**
   - [26]*TrieNode instead of map
   - Fixed size, no map overhead

2. **Bit packing:**
   - Pack isEnd as bit in pointer (alignment)
   - Only 2^48 address space needed

3. **Lazy initialization:**
   - Don't create children map until needed
   - Use nil to indicate no children

4. **Alphabet compression:**
   - Map rare characters to smaller set
   - Store mapping table

5. **Suffix tree:**
   - For static datasets
   - More complex but very space efficient

---

### 23. Persistence
**Question:** How would you save/load trie to disk?

**Answer:**

**Serialization approaches:**

1. **JSON/Text:**
```go
type SerialNode struct {
    Char     rune
    IsEnd    bool
    Children []SerialNode
}
```

2. **Binary format:**
- Pre-order traversal
- Write: char | isEnd flag | child count
- More compact than JSON

3. **Memory mapping:**
- Fixed-size nodes in array
- Pointers as indices
- Can mmap() for fast loading

---

## Practice Drill

**Warm-up (5 minutes each):**
1. Draw a trie after inserting: "car", "cat", "dog"
2. Explain why root is a dummy node
3. Walk through searching for "ca" vs "cat"

**Medium (10-15 minutes each):**
4. Implement Insert and Search on whiteboard
5. Implement FindAllWithPrefix for autocomplete
6. Draw and explain Delete("card") from trie: "car", "card", "cat"

**Advanced (20+ minutes):**
7. Implement word search with '.' wildcard
8. Design autocomplete with frequency ranking
9. Implement compressed trie (radix tree)

---

## Key Takeaways for Interviews

1. **Root is dummy** - Doesn't represent any character
2. **Path spells word** - Map keys form the word, not node values
3. **isEnd flag essential** - Distinguishes "car" from "card"
4. **Use rune not byte** - Supports Unicode properly
5. **Initialize children map** - Don't forget make()
6. **Prefix power** - O(m) prefix queries vs O(n) for other structures
7. **Space trade-off** - More pointers but shared prefixes
8. **Autocomplete king** - Perfect for typeahead/suggestions

**Interview Tips:**
- Start by drawing a small example trie
- Explain root as dummy/sentinel node
- Clarify alphabet (a-z? Unicode? case-sensitive?)
- Discuss space optimization for given alphabet
- Test with: empty trie, single char, shared prefixes
- Know real-world uses: autocomplete, IP routing, spell check

**Common interview problems:**
- Implement basic operations (Insert/Search/StartsWith)
- Word Search II (Boggle)
- Add and Search with wildcards
- Replace words with roots
- Autocomplete with ranking

Good luck!
